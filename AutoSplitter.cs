using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using System.Windows.Forms;
using System.Drawing;

namespace hksplitmaker
{
    class LineData
    {
        private static int INDEX = 0;
        private Panel line;
        private TextBox name;
        private ComboBox splitId;
        private Button delBtn;
        private Button addUpBtn;
        private Button addDownBtn;
        private Panel parent;
        public LineData(Panel parent)
        {
            int i = INDEX++;
            name = new TextBox();
            name.Location = new Point(17, 1);
            name.Name = "nameBox" + i.ToString();
            name.Size = new Size(200, 23);
            name.TabIndex = 0;
            splitId = new ComboBox();
            splitId.FormattingEnabled = true;
            splitId.Location = new Point(222, 0);
            splitId.Name = "splitIdComboBox" + i.ToString();
            splitId.Size = new Size(260, 25);
            splitId.TabIndex = 1;
            AutoSplitter.Instance().InitComboBox(splitId);
            delBtn = new Button();
            delBtn.Location = new Point(487, 0);
            delBtn.Name = "delBtn" + i.ToString();
            delBtn.Size = new Size(35, 25);
            delBtn.TabIndex = 2;
            delBtn.Text = "✘";
            delBtn.UseVisualStyleBackColor = true;
            delBtn.Click += DelBtn_Click;
            addUpBtn = new Button();
            addUpBtn.Location = new Point(525, 0);
            addUpBtn.Name = "addUpBtn" + i.ToString();
            addUpBtn.Size = new Size(35, 25);
            addUpBtn.TabIndex = 3;
            addUpBtn.Text = "↑+";
            addUpBtn.UseVisualStyleBackColor = true;
            addDownBtn = new Button();
            addDownBtn.Location = new Point(563, 0);
            addDownBtn.Name = "addDownBtn" + i.ToString();
            addDownBtn.Size = new Size(35, 25);
            addDownBtn.TabIndex = 4;
            addDownBtn.Text = "↓+";
            addDownBtn.UseVisualStyleBackColor = true;
            line = new Panel();
            line.TabIndex = i + 100;
            line.Controls.Add(name);
            line.Controls.Add(splitId);
            line.Controls.Add(delBtn);
            line.Controls.Add(addUpBtn);
            line.Controls.Add(addDownBtn);
            line.Location = new Point(0, 45 + i * 32);
            line.Name = "autoSpliterLine" + i.ToString();
            line.Size = new Size(parent.Width, 28);
            this.parent = parent;
            this.parent.Controls.Add(line);
            line.ResumeLayout(false);
            line.PerformLayout();
        }

        private void DelBtn_Click(object sender, EventArgs e)
        {
            if (INDEX > 1)
            {
                this.parent.Controls.RemoveAt(--INDEX + 100);
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

        public static AutoSplitter Instance()
        {
            if (instance == null)
            {
                instance = new AutoSplitter();
                instance.Init();
            }
            return instance;
        }

        private void initSplitsSearchDict(string content)
        {
            for (int i = 0; i < content.Length; i++)
            {
                for (int j = 1; j <= content.Length-i; j++)
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
            IEnumerable<string> lines = File.ReadLines("../../../splits.txt");
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
                    } else
                    {
                        throw new Exception("splits.txt文件格式错误");
                    }
                } else
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
        }
    }
}
