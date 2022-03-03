using System.Xml;

namespace hksplitmaker
{
    class SplitFile
    {
        public static void WriteFile(string filename)
        {
            XmlDocument xmlDoc = new XmlDocument();
            XmlDeclaration xmlDec = xmlDoc.CreateXmlDeclaration("1.0", "utf-8", null);
            xmlDoc.AppendChild(xmlDec);

            XmlElement rootElement = xmlDoc.CreateElement("Run");
            rootElement.SetAttribute("version", "1.7.0");

            XmlElement gameIcon = xmlDoc.CreateElement("GameIcon");
            gameIcon.InnerText = "";
            rootElement.AppendChild(gameIcon);

            XmlElement gameName = xmlDoc.CreateElement("GameName");
            gameName.InnerText = "Hollow Knight";
            rootElement.AppendChild(gameName);

            XmlElement categoryName = xmlDoc.CreateElement("CategoryName");
            categoryName.InnerText = "";
            rootElement.AppendChild(categoryName);

            XmlElement metaData = xmlDoc.CreateElement("Metadata");
            XmlElement metaDataRun = xmlDoc.CreateElement("Run");
            metaDataRun.SetAttribute("id", "");
            metaData.AppendChild(metaDataRun);
            XmlElement metaDataPlatform = xmlDoc.CreateElement("Platform");
            metaDataPlatform.SetAttribute("usesEmulator", "False");
            metaData.AppendChild(metaDataPlatform);
            XmlElement metaDataVariables = xmlDoc.CreateElement("Variables");
            metaDataVariables.InnerText = "";
            metaData.AppendChild(metaDataVariables);
            rootElement.AppendChild(metaData);

            XmlElement offset = xmlDoc.CreateElement("Offset");
            offset.InnerText = "00:00:00";
            rootElement.AppendChild(offset);

            XmlElement attemptCount = xmlDoc.CreateElement("AttemptCount");
            attemptCount.InnerText = "0";
            rootElement.AppendChild(attemptCount);

            XmlElement attemptHistory = xmlDoc.CreateElement("AttemptHistory");
            attemptHistory.InnerText = "";
            rootElement.AppendChild(attemptHistory);

            XmlElement segments = xmlDoc.CreateElement("Segments");
            foreach (LineData data in LineData.All)
            {
                XmlElement segment = xmlDoc.CreateElement("Segment");
                XmlElement segmentName = xmlDoc.CreateElement("Name");
                segmentName.InnerText = data.NameText;
                segment.AppendChild(segmentName);
                XmlElement segmentIcon = xmlDoc.CreateElement("Icon");
                segmentIcon.InnerText = "";
                segment.AppendChild(segmentIcon);
                XmlElement segmentSplitTimes = xmlDoc.CreateElement("SplitTimes");
                XmlElement segmentSplitTime = xmlDoc.CreateElement("SplitTime");
                segmentSplitTime.SetAttribute("name", "Personal Best");
                segmentSplitTime.InnerText = "";
                segmentSplitTimes.AppendChild(segmentSplitTime);
                segment.AppendChild(segmentSplitTimes);
                segments.AppendChild(segment);
            }
            XmlElement finalLineSegment = xmlDoc.CreateElement("Segment");
            XmlElement finalLineSegmentName = xmlDoc.CreateElement("Name");
            finalLineSegmentName.InnerText = FinalLineData.Instance.NameText;
            finalLineSegment.AppendChild(finalLineSegmentName);
            XmlElement finalLineSegmentIcon = xmlDoc.CreateElement("Icon");
            finalLineSegmentIcon.InnerText = "";
            finalLineSegment.AppendChild(finalLineSegmentIcon);
            XmlElement finalLineSegmentSplitTimes = xmlDoc.CreateElement("SplitTimes");
            XmlElement finalLineSegmentSplitTime = xmlDoc.CreateElement("SplitTime");
            finalLineSegmentSplitTime.SetAttribute("name", "Personal Best");
            finalLineSegmentSplitTime.InnerText = "";
            finalLineSegmentSplitTimes.AppendChild(finalLineSegmentSplitTime);
            finalLineSegment.AppendChild(finalLineSegmentSplitTimes);
            segments.AppendChild(finalLineSegment);
            rootElement.AppendChild(segments);

            XmlElement autoSplitterSettings = xmlDoc.CreateElement("AutoSplitterSettings");
            XmlElement ordered = xmlDoc.CreateElement("Ordered");
            ordered.InnerText = "True";
            autoSplitterSettings.AppendChild(ordered);
            XmlElement autosplitEndRuns = xmlDoc.CreateElement("AutosplitEndRuns");
            autosplitEndRuns.InnerText = FinalLineData.Instance.EndTrigger ? "False" : "True";
            autoSplitterSettings.AppendChild(autosplitEndRuns);
            XmlElement autosplitStartRuns = xmlDoc.CreateElement("AutosplitStartRuns");
            autosplitStartRuns.InnerText = "";
            autoSplitterSettings.AppendChild(autosplitStartRuns);
            XmlElement splits = xmlDoc.CreateElement("Splits");
            foreach (LineData data in LineData.All)
            {
                XmlElement split = xmlDoc.CreateElement("Split");
                split.InnerText = AutoSplitter.Instance.DescriptionToId(data.SplitIdText);
                splits.AppendChild(split);
            }
            if (!FinalLineData.Instance.EndTrigger)
            {
                XmlElement split = xmlDoc.CreateElement("Split");
                split.InnerText = AutoSplitter.Instance.DescriptionToId(FinalLineData.Instance.SplitIdText);
                splits.AppendChild(split);
            }
            autoSplitterSettings.AppendChild(splits);
            rootElement.AppendChild(autoSplitterSettings);

            xmlDoc.AppendChild(rootElement);
            xmlDoc.Save(filename);
        }
    }
}
