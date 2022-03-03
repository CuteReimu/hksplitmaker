using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace hksplitmaker
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
            AutoSplitter.Instance().InitComboBox(this.comboBox2);
            LineData.AddLine(this.panel1);
            LineData.AddLine(this.panel1);
            LineData.AddLine(this.panel1);
            FinalLineData.Init(this.panel1);
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
    }
}
