using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Text.RegularExpressions;
using System.Windows.Forms;

namespace hksplitmaker
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
            LineData.Parent = this.panel1;
            AutoSplitter.Instance.InitComboBox(this.comboBox2);
            LineData.AddLine();
            FinalLineData.Init(this.panel1);
            Categories.Instance.InitComboBox(this.comboBox1, false);
        }

        private void button2_Click(object sender, EventArgs e)
        {
            MessageBox.Show("1. 你可以选择一个已有的模板或者打开一个lss文件，也可以自己从头开始编辑\n" +
    "2. 点击右侧的加号和删除可以增加或者删除一行\n" +
    "3. 第一列的文本框里填写自己想要显示的片段名称，第二列的下拉框用来选择在游戏中会自动分割的触发事件。\n" +
    "4. 如果是一个全关速通或者万神殿某一门的速通，不需要勾选【自动开始】\n" +
    "5. 最后一行的复选框的意思是：如果你想要以游戏结束（打出任何一个结局）为计时器结束的标志，那么请勾上这个复选框；如果你想要以并非游戏结束的一个事件作为计时器结束的标志，那么请不要勾选这个复选框。\n" +
    "6. 全部设置好以后，点击下方的【另存为】按钮，即可保存成Splits文件。\n" +
    "7. 打开LiveSplit -> 右键 -> Open Splits -> From File... ，选择刚刚保存的文件即可。\n\n" +
    "版权所有 © CuteReimu 奇葩の灵梦\n" +
    "本程序的Github项目地址是：https://github.com/CuteReimu/hksplitmaker\n\n" +
    "本程序的所有非代码部分（图标和模板）全部都来自：https://hksplitmaker.com/\n" +
    "该项目的Github地址是：https://github.com/slaurent22/hk-split-maker", "帮助", MessageBoxButtons.OK, MessageBoxIcon.Information);
        }

        private void checkBox1_CheckedChanged(object sender, EventArgs e)
        {
            this.comboBox2.Enabled = ((CheckBox)sender).Checked;
        }

        private void button3_Click(object sender, EventArgs e)
        {
            this.saveFileDialog1.ShowDialog();
        }

        private void saveFileDialog1_FileOk(object sender, CancelEventArgs e)
        {
            SplitFile.WriteFile(this.saveFileDialog1.FileName);
        }

        private string categoryCurrent;

        private void comboBox1_SelectedIndexChanged(object sender, EventArgs e)
        {
            string category = comboBox1.Text;
            if (category == null || category.Length == 0 || category == categoryCurrent)
            {
                return;
            }
            categoryCurrent = category;
            CategoryData data = Categories.Instance[category];
            int count = data.SplitIds.Count;
            if (!data.EndTriggeringAutosplit)
            {
                count++;
            }
            LineData.ResetLines(count - 1);
            Regex re = new Regex("{.*?}|\\[[0-9DU, ]*]");
            IDictionary<string, int> nameIndexCache = new Dictionary<string, int>();
            Func<string, string> dropBrackets = new Func<string, string>((string s) =>
            {
                int idx = s.LastIndexOf('(');
                if (idx > 0)
                {
                    return s.Substring(0, idx);
                }
                return s;
            });
            Func<string, string, string> getNameFunc = new Func<string, string, string>((string splitId, string description) =>
            {
                string name = "";
                if (data.Names != null && data.Names.ContainsKey(splitId))
                {
                    Object names = data.Names[splitId];
                    if (names is string)
                    {
                        name = ((string)names).Replace("%s", dropBrackets(description));
                    }
                    else
                    {
                        object[] namearr = (object[])names;
                        if (!nameIndexCache.ContainsKey(splitId))
                        {
                            nameIndexCache[splitId] = 0;
                        }
                        name = ((string)namearr[nameIndexCache[splitId]]).Replace("%s", dropBrackets(description));
                        nameIndexCache[splitId]++;
                    }
                }
                else
                {
                    name = dropBrackets(description);
                }
                return re.Replace(name, "").Trim();
            });
            string startTrigger = AutoSplitter.Instance.IdToDescription(data.StartTriggeringAutosplit);
            if (startTrigger != null)
            {
                this.checkBox1.Checked = true;
                this.comboBox2.Text = startTrigger;
                this.comboBox2.Enabled = true;
            }
            else
            {
                this.checkBox1.Checked = false;
                this.comboBox2.Enabled = false;
            }
            if (data.EndTriggeringAutosplit)
            {
                for (int i = 0; i < data.SplitIds.Count; i++)
                {
                    string splitId = data.SplitIds[i].Trim('-');
                    string description = AutoSplitter.Instance.IdToDescription(splitId);
                    if (i < data.SplitIds.Count - 1)
                    {
                        LineData.All[i].SplitIdText = description;
                        LineData.All[i].NameText = getNameFunc(splitId, description);
                    }
                    else
                    {
                        FinalLineData.Instance.SetEndTrigger(false, description);
                    }
                }
            }
            else
            {
                for (int i = 0; i < data.SplitIds.Count; i++)
                {
                    string splitId = data.SplitIds[i].Trim('-');
                    string description = AutoSplitter.Instance.IdToDescription(splitId);
                    LineData.All[i].SplitIdText = description;
                    LineData.All[i].NameText = getNameFunc(splitId, description);
                }
                string text = "空洞骑士";
                if (data.EndingSplit.Icon == "RadianceBoss")
                {
                    text = data.EndingSplit.Name == "Absolute Radiance" ? "无上辐光" : "辐光";
                }
                FinalLineData.Instance.SetEndTrigger(true, text);
                FinalLineData.Instance.NameText = getNameFunc(data.EndingSplit.Icon, data.EndingSplit.Name);
            }
            this.checkBox2.Enabled = false;
            this.checkBox2.Checked = false;
            FinalLineData.UpdateLocation();
        }
    }
}
