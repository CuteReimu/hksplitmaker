using System;
using System.Collections.Generic;
using System.Drawing;
using System.Text.RegularExpressions;
using System.Windows.Forms;

namespace hksplitmaker
{
    class Trie
    {
        private static ISet<char> symbols = new HashSet<char>(new char[] { ' ', '(', ')', '[', ']', '-', '{', '}', '%', '\'', '"' });
        private class TrieNode
        {
            public IDictionary<char, TrieNode> child;
            public string value;
        }

        private TrieNode root = new TrieNode();

        public bool PutIfAbsent(string key, string value)
        {
            TrieNode node = this.root;
            foreach (char c in key.ToLower())
            {
                if (node.child == null)
                {
                    node.child = new Dictionary<char, TrieNode>();
                }
                TrieNode n = null;
                if (node.child.TryGetValue(c, out n))
                {
                    node = n;
                } else
                {
                    TrieNode newNode = new TrieNode();
                    node.child[c] = newNode;
                    node = newNode;
                }
            }
            if (node.value != null)
            {
                return false;
            }
            node.value = value;
            return true;
        }

        private string GetLongest(string s, out string key2)
        {
            s = s.ToLower();
            TrieNode node = this.root, node2 = null;
            string key = "";
            key2 = null;
            for (int i = 0; i < s.Length; i++)
            {
                if (node.child != null)
                {
                    char c = s[i];
                    if (node.child.ContainsKey(c))
                    {
                        key += c;
                        node = node.child[c];
                        if (node.value != null && (i+1 >= s.Length || symbols.Contains(s[i+1]))) {
                            node2 = node;
                            key2 = key;
                        }
                        continue;
                    }
                }
                break;
            }
            return node2 != null ? node2.value : null;
        }

        public string ReplaceAll(string s)
        {
            string s2 = "";
            while (s.Length > 0)
            {
                if (!(s2.Length == 0 || symbols.Contains(s2[s2.Length - 1])))
                {
                    s2 += s[0];
                    s = s.Substring(1);
                    continue;
                }
                string key;
                string value = GetLongest(s, out key);
                if (key != null)
                {
                    s2 += value;
                    s = s.Substring(key.Length);
                } else
                {
                    s2 += s[0];
                    s = s.Substring(1);
                }
            }
            return s2;
        }
    }

    class Translator
    {
        private static Regex regexpSpace = new Regex("(?<![()\\[\\]{}%'\"A - Za - z]) (?![()\\[\\]{}%'\"A-Za-z])");

        private static Translator instance;

        public static Translator Instance
        {
            get
            {
                if (instance == null)
                {
                    instance = new Translator();
                    instance.Init();
                }
                return instance;
            }
        }

        private Trie translateDict = new Trie();

        private Translator() { }

        private void Init()
        {
            foreach (string line1 in Resource.ResourceManager.GetString("translate.tsv").Split("\n"))
            {
                string line = line1.Trim();
                if (line.Length > 0)
                {
                    string[] arr = line.Split("\t");
                    string key = arr[0];
                    string val = arr.Length >= 2 ? arr[1] : "";
                    if (!translateDict.PutIfAbsent(key, val))
                    {
                        MessageBox.Show("出现重复数据：" + line, "警告", MessageBoxButtons.OK, MessageBoxIcon.Warning);
                    }
                }
            }
        }

        public string Translate(string s)
        {
            return regexpSpace.Replace(translateDict.ReplaceAll(s), "");
        }
    }
}
