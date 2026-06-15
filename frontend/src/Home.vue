<template>
  <el-alert title="如果想为本项目做贡献，请前往本项目的仓库地址： https://github.com/CuteReimu/hksplitmaker" effect="dark" close-text="前往" @close="openGithub" style="max-width: 960px"></el-alert>
  <div>
    <el-select v-model="currentTemplate" filterable placeholder="你可以选择现有模板" style="width: 400px" @change="selectTemplate">
        <el-option v-for="item in templates" :key="item.value" :label="item.label" :value="item.value"></el-option>
    </el-select>
    <el-select v-model="otherTemplate" filterable placeholder="你也可以选择其他玩家友情提供的模板" style="width: 400px; margin-left: 50px" @change="selectOtherTemplate">
      <el-option v-for="item in otherTemplates" :key="item.value" :label="item.label" :value="item.value"></el-option>
    </el-select>
  </div>
  <el-upload drag accept=".lss" :auto-upload="false" :show-file-list="false" :on-change="handleChange">
    <el-icon class="el-icon--upload"><upload-filled style="width: 80px;"></upload-filled></el-icon>
    <div class="el-upload__text">
      你也可以将文件拖拽到这里或者 <em>点击上传</em> 只支持 *.lss 文件
    </div>
  </el-upload>
  <div style="display: flex; gap: 8px;">
    <el-button @click="fillIcons">一键填充所有未填充的图标</el-button>
    <el-button @click="resetIcons">一键清空所有图标</el-button>
    <el-button @click="downloadIcons">下载所有图标</el-button>
    <el-text style="margin: 0 10px;">Auto Splitter版本：3.2.6.0</el-text>
    <el-button @click="fixLiveSplit" :loading="fixingLiveSplit">更新LiveSplit</el-button>
    <el-button @click="openHelp">查看帮助</el-button>
  </div>
  <el-table :data="tableData" max-width="960px">
    <el-table-column label="图标" width="60px">
      <template #default="scope">
        <el-image v-if="scope.$index>0 && scope.row.icon.length>0 && enableTriggering(scope.$index)" style="width: 25px; height: 25px" :src="scope.row.icon" fit="contain"></el-image>
      </template>
    </el-table-column>
    <el-table-column label="节点名称">
      <template #default="scope">
        <el-input v-model="scope.row.name" :disabled="scope.$index==0" placeholder="节点名称" style="width: 300px"></el-input>
      </template>
    </el-table-column>
    <el-table-column label="触发事件">
      <template #default="scope">
        <el-select-v2 v-if="enableTriggering(scope.$index)" v-model="scope.row.event" :options="options"
                   @change="onEventChange(scope.$index)" filterable placeholder="触发事件" style="width: 300px" />
      </template>
    </el-table-column>
    <el-table-column label="操作" :width="220">
      <template #default="scope">
        <el-button v-if="scope.$index>0" :icon="Plus" circle @click="addLine(scope.$index)"></el-button>
        <el-button :disabled="tableData.length<=3" v-if="scope.$index>0 && scope.$index<tableData.length-1" :icon="Minus" circle @click="removeLine(scope.$index)"></el-button>
        <el-button :disabled="scope.$index<=1" v-if="scope.$index>0 && scope.$index<tableData.length-1" :icon="Top" circle @click="swapLine(scope.$index-1, scope.$index)"></el-button>
        <el-button :disabled="scope.$index>=tableData.length-2" v-if="scope.$index>0 && scope.$index<tableData.length-1" :icon="Bottom" @click="swapLine(scope.$index, scope.$index+1)" circle></el-button>
        <el-checkbox v-if="scope.$index==0" v-model="startTriggering" style="margin-left:10px">自定义开始节点</el-checkbox>
        <el-checkbox v-if="scope.$index==tableData.length-1" v-model="endTriggering" style="margin-left:10px">不是游戏结束节点</el-checkbox>
      </template>
    </el-table-column>
  </el-table>
  <div>
    <el-button type="primary" @click="submit" style="align-self: flex-start;" :disabled="disableSubmit">另存为</el-button>
    <el-checkbox v-model="includeTimeRecords" size="large" style="margin-left: 20px">保留*.lss文件中原本的时间记录（如果有）</el-checkbox>
  </div>
</template>

<script setup lang="ts">
import {
  ElAlert,
  ElSelect,
  ElSelectV2,
  ElOption,
  ElUpload,
  ElButton,
  ElTable,
  ElTableColumn,
  ElCheckbox,
  ElMessage,
  ElText,
  ElIcon,
  ElImage,
  ElInput,
  UploadFile,
} from 'element-plus';
import { Plus, Minus, Top, Bottom, UploadFilled } from '@element-plus/icons-vue';
import {ref, onMounted} from 'vue';
import { GetOptions, GetTemplates, LoadSplitFile, GetSplits, GetIcon, SaveSplitsFile, SaveIconsZip, FixLiveSplit, GetUserDefinedFiles, OnSelectUserDefinedFile } from '../wailsjs/go/main/App';
import {BrowserOpenURL, LogError, EventsOn, EventsEmit} from '../wailsjs/runtime';

interface Row {
  name: string;
  event: string;
  icon: string;
  other?: unknown[];
}

interface Option {
  value: string;
  label: string;
}

const includeTimeRecords = ref(true);
const disableSubmit = ref(false);
const options = ref<Option[]>([]);
const currentTemplate = ref('');
const otherTemplate = ref('');
const templates = ref<Option[]>([]);
const otherTemplates = ref<Option[]>([]);
const fixingLiveSplit = ref(false);
const tableData = ref<Row[]>([
  { name: '开始', event: 'StartNewGame', icon: '' },
  { name: '空洞骑士', event: 'HollowKnightBoss', icon: '' },
]);
const startTriggering = ref(false);
const endTriggering = ref(false);
const enableTriggering = (index: number) =>
  (index > 0 && index < tableData.value.length - 1) ||
  (index === 0 && startTriggering.value) ||
  index === tableData.value.length - 1;

onMounted(() => {
  GetOptions().then(res => {
    options.value = res;
  }).catch(e => {
    LogError(e);
  });
  GetTemplates().then(res => {
    templates.value = res;
  }).catch(e => {
    LogError(e);
  });
  GetUserDefinedFiles().then(res => {
    otherTemplates.value = res;
  }).catch(e => {
    LogError(e);
  });
  GetIcon(tableData.value[1].event).then(res => {
    tableData.value[1].icon = res;
  }).catch(e => {{
    LogError(e);
  }});
});

function addLine(index: number) {
  GetIcon('ManualSplit').then(res => {
    tableData.value.splice(index, 0, { name: '手动分割', event: 'ManualSplit', icon: res });
  }).catch(e => {
    LogError(e);
    tableData.value.splice(index, 0, { name: '手动分割', event: 'ManualSplit', icon: '' });
  })
}

function removeLine(index: number) {
  tableData.value.splice(index, 1);
}

function swapLine(index1: number, index2: number) {
  const temp = tableData.value[index1];
  tableData.value[index1] = tableData.value[index2];
  tableData.value[index2] = temp;
}

function submit() {
  disableSubmit.value = true;
  SaveSplitsFile(
    tableData.value as any,
    includeTimeRecords.value,
    startTriggering.value,
    endTriggering.value,
  ).catch(e => {
    LogError(e);
    ElMessage({ message: '导出失败', type: 'error', plain: true });
  }).finally(() => {
    disableSubmit.value = false;
  });
}

function downloadIcons() {
  disableSubmit.value = true;
  SaveIconsZip().catch(e => {
    LogError(e);
    ElMessage({ message: '导出失败', type: 'error', plain: true });
  }).finally(() => {
    disableSubmit.value = false;
  });
}

function handleChange(file: UploadFile) {
  if (!file?.raw) return;
  file.raw.text().then(text => {
    LoadSplitFile(text).then(newData => {
      tableData.value = newData.splits as Row[];
      startTriggering.value = newData.startTriggering;
      endTriggering.value = newData.endTriggering;
    }).catch(e => {
      LogError(e);
      ElMessage({ message: String(e), type: 'error', plain: true });
    })
  }).catch(e => {
    LogError(e);
    ElMessage({ message: String(e), type: 'error', plain: true });
  });
}

const coloEnd = ["BronzeEnd", "SilverEnd", "GoldEnd"];

function selectTemplate(value: string) {
  otherTemplate.value = '';
  GetSplits(value).then(res => {
    startTriggering.value = res.startTriggering;
    endTriggering.value = res.endTriggering;
    tableData.value = res.splits;
    if (res.splits.some(v => coloEnd.includes(v.event))) {
      EventsEmit("onSelectColo");
    }
  }).catch(e => {
    LogError(e);
    ElMessage({ message: String(e), type: 'error', plain: true });
  });
}

function selectOtherTemplate(value: string) {
  currentTemplate.value = '';
  OnSelectUserDefinedFile(value).then(res => {
    startTriggering.value = res.startTriggering;
    endTriggering.value = res.endTriggering;
    tableData.value = res.splits;
  }).catch(e => {
    LogError(e);
    ElMessage({ message: String(e), type: 'error', plain: true });
  });
}

function openGithub() {
  BrowserOpenURL('https://github.com/CuteReimu/hksplitmaker');
}

function openHelp() {
  BrowserOpenURL('https://cutereimu.cn/daily/hollowknight/hksplitmaker-faq.html');
}

function onEventChange(idx: number) {
  if (idx === 0) {
    return;
  }
  const eventValue = tableData.value[idx].event;
  const opt = options.value.find(o => o.value === eventValue);
  if (opt) {
    const pos = opt.label.indexOf('（');
    tableData.value[idx].name = pos === -1 ? opt.label : opt.label.slice(0, pos);
  }
  GetIcon(tableData.value[idx].event).then(icon => {
    tableData.value[idx].icon = icon;
  }).catch(e => {
    LogError(e);
    ElMessage({ message: String(e), type: 'error', plain: true });
  });
}

function fillIcons() {
  const p = [];
  for (const idx in tableData.value) {
    const row = tableData.value[idx];
    if (row.icon.length === 0) {
      p.push(GetIcon(row.event).then(icon => {
        row.icon = icon;
      }));
    }
  }
  Promise.all(p).catch(e => {
    LogError(e);
    ElMessage({ message: String(e), type: 'error', plain: true });
  })
}

function resetIcons() {
  for (const idx in tableData.value) {
    tableData.value[idx].icon = '';
  }
}

function fixLiveSplit() {
  fixingLiveSplit.value = true;
  FixLiveSplit().finally(() => {
    fixingLiveSplit.value = false;
  })
}

EventsOn("ElMessage", (type, message) => {
  ElMessage({ message, type, plain: true });
});
</script>
