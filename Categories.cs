using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.IO;
using System.Runtime.Serialization;
using System.Runtime.Serialization.Json;
using System.Text;
using System.Windows.Forms;

namespace hksplitmaker
{
    [DataContract]
    class CategoryData
    {
        [DataMember(Name = "categoryName", Order = 0)]
        public string CategoryName;

        [DataMember(Name = "splitIds", Order = 1)]
        public List<string> SplitIds = new List<string>();

        [DataMember(Name = "ordered", Order = 2)]
        public bool Ordered;

        [DataMember(Name = "startTriggeringAutosplit", Order = 3)]
        public string StartTriggeringAutosplit;

        [DataMember(Name = "endTriggeringAutosplit", Order = 4)]
        public bool EndTriggeringAutosplit;

        [DataMember(Name = "names", Order = 5)]
        public Dictionary<string, Object> Names = new Dictionary<string, Object>();

        [DataMember(Name = "icons", Order = 6)]
        public Dictionary<string, Object> Icons = new Dictionary<string, object>();

        [DataMember(Name = "endingSplit", Order = 7)]
        public CategoryEndSplit EndingSplit = new CategoryEndSplit();

        [DataMember(Name = "gameName", Order = 8)]
        public string GameName;
    }

    [DataContract]
    class CategoryEndSplit
    {
        [DataMember(Name = "name", Order = 0)]
        public string Name;

        [DataMember(Name = "icon", Order = 1)]
        public string Icon;
    }

    [DataContract]
    class CategoryInfo
    {
        [DataMember(Name = "fileName", Order = 0)]
        public string FileName;

        [DataMember(Name = "displayName", Order = 1)]
        public string DisplayName;
    }

    class Categories
    {
        private static Categories instance;

        public static Categories Instance
        {
            get
            {
                if (instance == null)
                {
                    instance = new Categories();
                    instance.Init();
                }
                return instance;
            }
        }

        private Dictionary<string, CategoryData> cache = new Dictionary<string, CategoryData>();

        public CategoryData this[string name] { get { return cache[name]; } }

        private Categories() { }

        private void Init()
        {
            Dictionary<string, JArray> categoryInfos = JsonConvert.DeserializeObject<Dictionary<string, JArray>>(Resource.category_directory_json);
            foreach (JArray infos in categoryInfos.Values)
            {
                foreach (JToken info in infos)
                {
                    CategoryData data = null;
                    string fileName = info["fileName"].ToString().Replace('-', '_') + ".json";
                    using (MemoryStream ms = new MemoryStream(Encoding.UTF8.GetBytes(Resource.ResourceManager.GetString(fileName))))
                    {
                        DataContractJsonSerializer jsonSerializer = new DataContractJsonSerializer(typeof(CategoryData));
                        data = (CategoryData)jsonSerializer.ReadObject(ms);
                    }
                    int count = data.SplitIds.Count;
                    if (!data.EndTriggeringAutosplit)
                    {
                        if (data.EndingSplit == null)
                        {
                            continue;
                        }
                        if (data.EndingSplit.Icon != "HollowKnightBoss" && data.EndingSplit.Icon != "RadianceBoss")
                        { // 暂时不支持
                            continue;
                        }
                        count++;
                    }
                    if (new Func<bool>(() =>
                    {
                        foreach (string splitId in data.SplitIds)
                            if (splitId.Contains("%"))
                                return true;
                        return false;
                    })())
                    { // 暂时不支持
                        continue;
                    }
                    if (count < 2 || !data.Ordered)
                    { // 暂时不支持
                        continue;
                    }
                    cache[info["displayName"].ToString()] = data;
                }
            }
        }

        public void InitComboBox(ComboBox b, bool selectFirst = true)
        {
            foreach (string s in cache.Keys)
            {
                b.Items.Add(s);
            }
            if (selectFirst)
            {
                b.SelectedIndex = 0;
            }
        }
    }
}
