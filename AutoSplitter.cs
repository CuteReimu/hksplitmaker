using System;
using System.Collections.Generic;
using System.Drawing;
using System.Text.RegularExpressions;
using System.Windows.Forms;

namespace hksplitmaker
{
    class FinalLineData
    {
        public static FinalLineData Instance { get; private set; }
        private readonly Panel line;
        private readonly TextBox name;
        private readonly ComboBox splitId, splitId2;
        private readonly CheckBox endTrigger;

        public string NameText { get { return name.Text; } set { name.Text = value; } }

        public string SplitIdText { get { return splitId.Text; } }

        public bool EndTrigger { get { return endTrigger.Checked; } }

        public void SetEndTrigger(bool check, string endTrigger)
        {
            this.endTrigger.Checked = check;
            if (check)
            {
                splitId2.Text = endTrigger;
            }
            else
            {
                splitId.Text = endTrigger;
            }
        }

        private FinalLineData()
        {
            name = new TextBox();
            name.Location = new Point(17, 1);
            name.Name = "finalLineNameBox";
            name.Size = new Size(200, 23);
            name.TabIndex = 0;
            splitId = new ComboBox();
            splitId.FormattingEnabled = true;
            splitId.Location = new Point(222, 0);
            splitId.Name = "finalLineSplitIdComboBox";
            splitId.Size = new Size(252, 25);
            splitId.TabIndex = 1;
            splitId.Visible = false;
            AutoSplitter.Instance.InitComboBox(splitId);
            splitId2 = new ComboBox();
            splitId2.FormattingEnabled = true;
            splitId2.Location = new Point(222, 0);
            splitId2.Name = "finalLineSplitIdComboBox2";
            splitId2.Size = new Size(252, 25);
            splitId2.TabIndex = 2;
            splitId2.Items.AddRange(new object[] { "空洞骑士", "辐光", "无上辐光" });
            splitId2.SelectedIndex = 0;
            endTrigger = new CheckBox();
            endTrigger.Location = new Point(479, 3);
            endTrigger.Name = "finalLineCheckBox";
            endTrigger.Size = new Size(106, 21);
            endTrigger.TabIndex = 0;
            endTrigger.Text = "游戏结束停止";
            endTrigger.UseVisualStyleBackColor = true;
            endTrigger.Checked = true;
            endTrigger.CheckedChanged += EndTrigger_CheckedChanged;
            line = new Panel();
            line.TabIndex = 3;
            line.Controls.Add(name);
            line.Controls.Add(splitId);
            line.Controls.Add(endTrigger);
            line.Controls.Add(splitId2);
            line.Location = new Point(0, 45 + LineData.Count * 32);
            line.Name = "finalLineAutoSpliterLine";
            line.Size = new Size(585, 28);
            LineData.Parent.Controls.Add(line);
            line.ResumeLayout(false);
            line.PerformLayout();
        }

        private void EndTrigger_CheckedChanged(object sender, EventArgs e)
        {
            this.splitId.Visible = !this.endTrigger.Checked;
            this.splitId2.Visible = this.endTrigger.Checked;
        }

        public static void Init(Panel parent) { Instance = new FinalLineData(); }

        public static void UpdateLocation() { Instance.line.Location = new Point(0, 45 + LineData.Count * 32); }
    }
    class LineData
    {
        private static IList<LineData> lineDataList = new List<LineData>();
        private static int INDEX = 0;
        private readonly Panel line;
        private readonly TextBox name;
        private readonly ComboBox splitId;
        private readonly Button delBtn;
        private readonly Button addUpBtn;
        private readonly Button addDownBtn;
        public static Panel Parent { set; get; }
        private readonly int index;
        public static LineData AddLine()
        {
            LineData lineData = new LineData(INDEX++);
            lineDataList.Add(lineData);
            return lineData;
        }

        public static void RemoveLine(int index)
        {
            for (int i = index + 1; i < lineDataList.Count; i++)
            {
                lineDataList[i - 1].name.Text = lineDataList[i].name.Text;
                lineDataList[i - 1].splitId.Text = lineDataList[i].splitId.Text;
            }
            lineDataList[lineDataList.Count - 1].line.Dispose();
            lineDataList.RemoveAt(lineDataList.Count - 1);
            INDEX--;
            FinalLineData.UpdateLocation();
        }

        public static int Count { get { return lineDataList.Count; } }

        public string NameText { get { return name.Text; } set { name.Text = value; } }

        public string SplitIdText { get { return splitId.Text; } set { splitId.Text = value; } }

        private LineData(int index)
        {
            this.index = index;
            name = new TextBox();
            name.Location = new Point(17, 1);
            name.Name = "nameBox" + index.ToString();
            name.Size = new Size(200, 23);
            name.TabIndex = 0;
            splitId = new ComboBox();
            splitId.FormattingEnabled = true;
            splitId.Location = new Point(222, 0);
            splitId.Name = "splitIdComboBox" + index.ToString();
            splitId.Size = new Size(252, 25);
            splitId.TabIndex = 1;
            AutoSplitter.Instance.InitComboBox(splitId);
            delBtn = new Button();
            delBtn.Location = new Point(478, 0);
            delBtn.Name = "delBtn" + index.ToString();
            delBtn.Size = new Size(35, 25);
            delBtn.TabIndex = 2;
            delBtn.Text = "✘";
            delBtn.UseVisualStyleBackColor = true;
            delBtn.Click += DelBtn_Click;
            addUpBtn = new Button();
            addUpBtn.Location = new Point(514, 0);
            addUpBtn.Name = "addUpBtn" + index.ToString();
            addUpBtn.Size = new Size(35, 25);
            addUpBtn.TabIndex = 3;
            addUpBtn.Text = "↑+";
            addUpBtn.UseVisualStyleBackColor = true;
            addUpBtn.Click += AddUpBtn_Click;
            addDownBtn = new Button();
            addDownBtn.Location = new Point(550, 0);
            addDownBtn.Name = "addDownBtn" + index.ToString();
            addDownBtn.Size = new Size(35, 25);
            addDownBtn.TabIndex = 4;
            addDownBtn.Text = "↓+";
            addDownBtn.UseVisualStyleBackColor = true;
            addDownBtn.Click += AddDownBtn_Click;
            line = new Panel();
            line.TabIndex = index + 100;
            line.Controls.Add(name);
            line.Controls.Add(splitId);
            line.Controls.Add(delBtn);
            line.Controls.Add(addUpBtn);
            line.Controls.Add(addDownBtn);
            line.Location = new Point(0, 45 + index * 32);
            line.Name = "autoSpliterLine" + index.ToString();
            line.Size = new Size(585, 28);
            Parent.Controls.Add(line);
            line.ResumeLayout(false);
            line.PerformLayout();
        }

        private void AddDownBtn_Click(object sender, EventArgs e)
        {
            LineData line = AddLine();
            for (int i = lineDataList.Count - 1; i > index + 1; i--)
            {
                lineDataList[i].name.Text = lineDataList[i - 1].name.Text;
                lineDataList[i].splitId.Text = lineDataList[i - 1].splitId.Text;
            }
            lineDataList[index + 1].name.Text = "";
            lineDataList[index + 1].splitId.SelectedIndex = 0;
            FinalLineData.UpdateLocation();
        }

        private void AddUpBtn_Click(object sender, EventArgs e)
        {
            LineData line = AddLine();
            for (int i = lineDataList.Count - 1; i >= index + 1; i--)
            {
                lineDataList[i].name.Text = lineDataList[i - 1].name.Text;
                lineDataList[i].splitId.Text = lineDataList[i - 1].splitId.Text;
            }
            lineDataList[index].name.Text = "";
            lineDataList[index].splitId.SelectedIndex = 0;
            FinalLineData.UpdateLocation();
        }

        private void DelBtn_Click(object sender, EventArgs e)
        {
            if (lineDataList.Count > 1)
            {
                RemoveLine(index);
            }
        }

        public static IList<LineData> All { get { return lineDataList; } }

        public static void ResetLines(int count)
        {
            for (int i = LineData.Count; i < count; i++)
            {
                LineData.AddLine();
            }
            for (int i = LineData.Count - 1; i >= count; i--)
            {
                LineData.RemoveLine(i);
            }
        }
    }

    class AutoSplitter
    {
        private static AutoSplitter instance;
        private IDictionary<string, string> idToDescription = new Dictionary<string, string>();
        private IDictionary<string, SplitData> descriptionToData = new SortedDictionary<string, SplitData>();
        private IDictionary<string, IList<string>> searchDict = new Dictionary<string, IList<string>>();
        private class SplitData
        {
            public string Description { get; }
            public string Tooltip { get; }
            public string Id { get; }
            public SplitData(string description, string tooltip, string id)
            {
                Description = description;
                Tooltip = tooltip;
                Id = id;
            }
        }

        public static AutoSplitter Instance
        {
            get
            {
                if (instance == null)
                {
                    instance = new AutoSplitter();
                    instance.Init();
                }
                return instance;
            }
        }

        private AutoSplitter() { }

        private void initSplitsSearchDict(string content)
        {
            for (int i = 0; i < content.Length; i++)
            {
                for (int j = 1; j <= content.Length - i; j++)
                {
                    string s = content.Substring(i, j);
                    if (!searchDict.ContainsKey(s))
                    {
                        searchDict[s] = new List<string>();
                    }
                    searchDict[s].Add(s);
                }
            }
        }

        private void Init()
        {
            string[] lines = Resource.splits_txt.Split("\n");
            bool isNameLine = false;
            string[] result = new string[0];
            Regex re = new Regex("\\[Description\\(\"(.*?)\"\\)\\s*,\\s*ToolTip\\(\"(.*?)\"\\)]");
            foreach (string l in lines)
            {
                string line = l.Trim().Trim(',');
                if (line.Length == 0)
                {
                    continue;
                }
                if (isNameLine)
                {
                    if (result.Length == 3)
                    {
                        string description = result[1];
                        descriptionToData[description] = new SplitData(description, result[2], line);
                        idToDescription[line] = description;
                        initSplitsSearchDict(description);
                        isNameLine = false;
                    }
                    else
                    {
                        throw new Exception("splits.txt文件格式错误");
                    }
                }
                else
                {
                    Match m = re.Match(line);
                    if (m == null)
                    {
                        throw new Exception("splits.txt文件格式错误");
                    }
                    result = new string[m.Groups.Count];
                    for (int i = 0; i < result.Length; i++)
                    {
                        result[i] = m.Groups[i].Value;
                    }
                    isNameLine = true;
                }
            }
        }

        public void InitComboBox(ComboBox b)
        {
            foreach (string s in descriptionToData.Keys)
            {
                b.Items.Add(s);
            }
            b.SelectedIndex = 0;
        }

        public string DescriptionToId(string description)
        {
            if (description == null)
            {
                return null;
            }
            SplitData data = descriptionToData[description];
            if (data == null)
            {
                return null;
            }
            return data.Id;
        }

        public string IdToDescription(string id)
        {
            if (id == null)
            {
                return null;
            }
            return idToDescription[id];
        }
    }
}
