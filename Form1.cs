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
            new LineData(this.panel1);
            new LineData(this.panel1);
            new LineData(this.panel1);
        }

    }
}
