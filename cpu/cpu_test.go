// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//Package cpu provides implementation of snap-plugin-collector-cpu plugin
package cpu

import (
	"os"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
)

type CPUInfoSuite struct {
	suite.Suite
	MockCPUInfo string
}

func (cis *CPUInfoSuite) SetupSuite() {
	cpuInfo = cis.MockCPUInfo
	loadMockCPUInfo(0)
}

func (cis *CPUInfoSuite) TearDownSuite() {
	removeMockCPUInfo()
}

func removeMockCPUInfo() {
	os.Remove(cpuInfo)
}

func TestGetStatsSuite(t *testing.T) {
	suite.Run(t, &CPUInfoSuite{MockCPUInfo: "MockCPUInfo"})
}

func mockNew() *Plugin {
	p := New()
	emptyCfg := make(map[string]ctypes.ConfigValue)
	p.init(emptyCfg)
	So(p, ShouldNotBeNil)
	So(p.snapMetricsNames, ShouldNotBeNil)
	return p
}

func loadMockCPUInfo(dataSetNumber int) {
	var content string
	if dataSetNumber == 0 {
		content = `cpu  23359837 6006716 1209900 402135131 129307 4 2156 0 0 0
			cpu0 3464284 998669 208226 49355234 57380 3 422 0 0 0
			cpu1 3501681 1012206 189642 49374240 11620 0 278 0 0 0
			intr 33594809 19 2 0 0 0 0 0 9 1 4 0 0 4 0 0 0 31 0 0`
	} else if dataSetNumber == 1 {
		content = `cpu  23472679 6048986 1215282 403105970 129312 4 2158 0 0 0
			cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0 0
			cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0
			intr 33594809 19 2 0 0 0 0 0 9 1 4 0 0 4 0 0 0 31 0 0`
	} else if dataSetNumber == 2 {
		content = `cpu  23472670 6049996 1215282 403105970 129312 4 2158 0 0 0
			cpu0 3480508 1005570 209105 49472590 57390 3 430 0 0 0
			cpu1 3516060 1019260 190410 49493310 11610 0 270 0 0 0`
	} else if dataSetNumber == 3 {
		content = `cpu  23472679 6048986 1215282 403105970 129312 4 2158 0 0 0
			cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0 0
			cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0`
	} else if dataSetNumber == 4 { //incorrect data set, should fail on parsing
		content = `cpu  23472679 6048986 1215282 403105970 129312 4 2158 0 0 0
			cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0
			cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0`
	} else if dataSetNumber == 5 { //incorrect data set, should fail on parsing
		content = ``
	} else if dataSetNumber == 6 { //incorrect data set, should fail on parsing
		content = `cpu  * # # 403105970 129312 4 2158 0 0 0
			cpu0 3480506 1005574 209103 49472588 57381 3 424 0 0 0
			cpu1 3516068 1019269 190413 49493320 11620 0 278 0 0 0
			intr 33594809 19 2 0 0 0 0 0 9 1 4 0 0 4 0 0 0 31 0 0`
	}

	cpuInfoContent := []byte(content)
	f, err := os.Create(cpuInfo)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.Write(cpuInfoContent)
}

func (cis *CPUInfoSuite) TestGetMetricTypes() {
	_ = plugin.ConfigType{}
	Convey("Given cpu info plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)
		Convey("When one wants to get list of available metrics", func() {
			mts, err := p.GetMetricTypes(plugin.ConfigType{})

			Convey("Then error should not be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then list of metrics is returned", func() {
				// Len mts = 24
				// cpuMetricsNumber = 3
				// len snapMetricsNames = 12
				So(len(mts), ShouldEqual, len(p.snapMetricsNames)*2)

				namespaces := []string{}
				for _, m := range mts {
					namespaces = append(namespaces, strings.Join(m.Namespace().Strings(), "/"))
				}

				So(namespaces, ShouldContain, "intel/procfs/cpu/*/user_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/nice_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/system_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/idle_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/iowait_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/irq_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/softirq_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/steal_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_nice_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/active_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/utilization_percentage")

				So(namespaces, ShouldContain, "intel/procfs/cpu/*/user_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/nice_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/system_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/idle_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/iowait_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/irq_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/softirq_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/steal_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_nice_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/active_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/utilization_jiffies")

			})
		})
	})
}

func (cis *CPUInfoSuite) TestCollectMetrics() {
	Convey("Given cpu plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)
		Convey("When one wants to get values for given metric types", func() {
			mTypes := []plugin.MetricType{
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(activeProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(utilizationProcStat, percentageRepresentationType))},

				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(activeProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, allCPU, getNamespaceMetricPart(utilizationProcStat, jiffiesRepresentationType))},

				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(activeProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(utilizationProcStat, percentageRepresentationType))},

				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(activeProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(utilizationProcStat, jiffiesRepresentationType))},

				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(activeProcStat, percentageRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, secondCPU, getNamespaceMetricPart(utilizationProcStat, percentageRepresentationType))},

				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(activeProcStat, jiffiesRepresentationType))},
				plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, firstCPU, getNamespaceMetricPart(utilizationProcStat, jiffiesRepresentationType))},
			}

			metrics, err := p.CollectMetrics(mTypes)

			Convey("Then no errors should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then proper metrics values are returned", func() {
				So(len(metrics), ShouldEqual, len(mTypes))
				for _, mt := range metrics {
					// Only jiffies metrics should not be nil
					for _, ns := range mt.Namespace_ {
						if strings.Contains(ns.Value, "jiffies") {
							So(mt.Data_, ShouldNotBeNil)
						}
					}
				}
			})
		})
	})
}

func (cis *CPUInfoSuite) TestgetCPUMetrics() {
	Convey("Given cpu plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)
		Convey("We want to check if metrics have proper value", func() {
			//get new data set from /proc/stat
			loadMockCPUInfo(0)
			errStats := getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
			So(errStats, ShouldBeNil)

			//all
			ns := core.NewNamespace(allCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))
			val, err := getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 23359837)
			_, ok := val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 6006716)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1209900)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 402135131)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 129307)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 4)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 2156)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//cpu0
			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3464284)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 998669)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 208226)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 49355234)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 57380)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 422)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//cpu1
			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3501681)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1012206)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 189642)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 49374240)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 11620)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 278)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldBeNil)

			//get new data set from /proc/stat
			loadMockCPUInfo(1)
			errStats = getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
			So(errStats, ShouldBeNil)

			//all
			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 23472679)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 6048986)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1215282)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 403105970)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 129312)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 4)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 2158)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//cpu0
			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3480506)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1005574)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 209103)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 49472588)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 57381)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 424)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//cpu1
			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(userProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 3516068)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(niceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 1019269)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(systemProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 190413)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(idleProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 49493320)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(iowaitProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 11620)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(irqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(softirqProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 278)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(stealProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestNiceProcStat, jiffiesRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//all percentage
			var prevAllSum, currAllSum float64
			prevAllSum = 23359837 + 6006716 + 1209900 + 402135131 + 129307 + 4 + 2156
			currAllSum = 23472679 + 6048986 + 1215282 + 403105970 + 129312 + 4 + 2158
			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(23472679-23359837)/(currAllSum-prevAllSum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(6048986-6006716)/(currAllSum-prevAllSum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(1215282-1209900)/(currAllSum-prevAllSum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(403105970-402135131)/(currAllSum-prevAllSum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(129312-129307)/(currAllSum-prevAllSum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)

			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(2158-2156)/(currAllSum-prevAllSum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//cpu0 percentage
			var prevCPU0Sum, currCPU0Sum float64
			prevCPU0Sum = 3464284 + 998669 + 208226 + 49355234 + 57380 + 3 + 422
			currCPU0Sum = 3480506 + 1005574 + 209103 + 49472588 + 57381 + 3 + 424
			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(3480506-3464284)/(currCPU0Sum-prevCPU0Sum))

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(1005574-998669)/(currCPU0Sum-prevCPU0Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(209103-208226)/(currCPU0Sum-prevCPU0Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(49472588-49355234)/(currCPU0Sum-prevCPU0Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(57381-57380)/(currCPU0Sum-prevCPU0Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(424-422)/(currCPU0Sum-prevCPU0Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			//cpu1 percentage
			var prevCPU1Sum, currCPU1Sum float64
			prevCPU1Sum = 3501681 + 1012206 + 189642 + 49374240 + 11620 + 278
			currCPU1Sum = 3516068 + 1019269 + 190413 + 49493320 + 11620 + 278
			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(3516068-3501681)/(currCPU1Sum-prevCPU1Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(1019269-1012206)/(currCPU1Sum-prevCPU1Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(190413-189642)/(currCPU1Sum-prevCPU1Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 100*(49493320-49374240)/(currCPU1Sum-prevCPU1Sum))
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
			val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
			So(err, ShouldBeNil)
			So(val, ShouldEqual, 0)
			_, ok = val.(float64)
			So(ok, ShouldBeTrue)

			Convey("We want to check if metric value is nil instead of negative in case of incorrect (decreasing) values in /proc/stat", func() {

				loadMockCPUInfo(1)
				errStats = getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
				So(errStats, ShouldBeNil)
				//get new data set to check percentage calculation for incorrect (decreasing) values in /proc/stat
				loadMockCPUInfo(2)
				errStats = getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
				So(errStats, ShouldBeNil)

				//all percentage
				// case 1: one value decreases, one increases (to keep the currAllSum bigger then the prevAllSum),
				// the rest does not change
				prevAllSum = 23472679 + 6048986 + 1215282 + 403105970 + 129312 + 4 + 2158
				currAllSum = 23472670 + 6049996 + 1215282 + 403105970 + 129312 + 4 + 2158
				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 100*(6049996-6048986)/(currAllSum-prevAllSum))
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(allCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				//cpu0 percentage
				// case 2: one value decreases, the rest increases (to check if all values are calculated correctly)
				prevCPU0Sum = 3480506 + 1005574 + 209103 + 49472588 + 57381 + 3 + 424
				currCPU0Sum = 3480508 + 1005570 + 209105 + 49472590 + 57390 + 3 + 430
				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 100*(3480508-3480506)/(currCPU0Sum-prevCPU0Sum))

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 100*(209105-209103)/(currCPU0Sum-prevCPU0Sum))
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 100*(49472590-49472588)/(currCPU0Sum-prevCPU0Sum))
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 100*(57390-57381)/(currCPU0Sum-prevCPU0Sum))
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 100*(430-424)/(currCPU0Sum-prevCPU0Sum))
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				ns = core.NewNamespace(firstCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldEqual, 0)
				_, ok = val.(float64)
				So(ok, ShouldBeTrue)

				//cpu1 percentage
				//case 3: all the values decrease
				prevCPU1Sum = 3516068 + 1019269 + 190413 + 49493320 + 11620 + 0 + 278
				currCPU1Sum = 3516060 + 1019260 + 190410 + 49493310 + 11610 + 0 + 270
				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(userProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(niceProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(systemProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(idleProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(iowaitProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(irqProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(softirqProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(stealProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)

				ns = core.NewNamespace(secondCPU, getNamespaceMetricPart(guestNiceProcStat, percentageRepresentationType))
				val, err = getMapValueByNamespace(p.stats[ns.Strings()[0]], ns.Strings()[1:])
				So(err, ShouldBeNil)
				So(val, ShouldBeNil)
			})

			Convey("We want to test getStats function with incorrect data sets", func() {
				loadMockCPUInfo(4)
				errStats = getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
				So(errStats, ShouldNotBeNil)

				loadMockCPUInfo(5)
				errStats = getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
				So(errStats, ShouldNotBeNil)

				loadMockCPUInfo(6)
				errStats = getStats(p.proc_path, p.stats, p.prevMetricsSum, p.cpuMetricsNumber, p.snapMetricsNames, p.procStatMetricsNames)
				So(errStats, ShouldNotBeNil)
			})
		})
	})
}

func (cis *CPUInfoSuite) TestgetInitialProcStatData() {
	Convey("Given cpu plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)
		Convey("We want to check initial reading of /proc/stat", func() {
			loadMockCPUInfo(1)
			_, _, err := getInitialProcStatData(p.proc_path)
			So(err, ShouldBeNil)
			loadMockCPUInfo(2)
			_, _, err = getInitialProcStatData(p.proc_path)
			So(err, ShouldBeNil)
			loadMockCPUInfo(3)
			_, _, err = getInitialProcStatData(p.proc_path)
			So(err, ShouldBeNil)
			loadMockCPUInfo(4)
			_, _, err = getInitialProcStatData(p.proc_path)
			So(err, ShouldNotBeNil)
			loadMockCPUInfo(5)
			_, _, err = getInitialProcStatData(p.proc_path)
			So(err, ShouldNotBeNil)
		})
	})
}

func (cis *CPUInfoSuite) TestgetMapFloatValueByNamespace() {
	Convey("Given cpu plugin initialized", cis.T(), func() {
		Convey("We want to check getting float value from nested map", func() {
			//create map to test
			subSubMap := make(map[string]interface{})
			subSubMap["subSubMap1"] = 56.0
			subSubMap["subSubMap2"] = "343.343"
			subSubMap["subSubMap3"] = '$'
			subSubMap["subSubMap4"] = 'a'

			subMap := make(map[string]interface{})
			subMap["subMap1"] = subSubMap

			mainMap := make(map[string]interface{})
			mainMap["mainMap1"] = subMap

			ns := core.NewNamespace("mainMap1", "subMap1", "subSubMap1")
			_, err := getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldBeNil)

			ns = core.NewNamespace("mainMap1", "subMap1", "subSubMap2")
			_, err = getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldNotBeNil)

			ns = core.NewNamespace("mainMap1", "subMap1", "subSubMap3")
			_, err = getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldNotBeNil)

			ns = core.NewNamespace("mainMap1", "subMap1", "subSubMap4")
			_, err = getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldNotBeNil)

			ns = core.NewNamespace("mainMap12", "subMap12", "subSubMap1")
			_, err = getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldNotBeNil)

			ns = core.NewNamespace("mainMap1", "subMap12", "subSubMap1")
			_, err = getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldNotBeNil)

			ns = core.NewNamespace("mainMap1", "subMap12", "subSubMap12")
			_, err = getMapFloatValueByNamespace(mainMap, ns.Strings())
			So(err, ShouldNotBeNil)

			var nss []string
			_, err = getMapFloatValueByNamespace(mainMap, nss)
			So(err, ShouldNotBeNil)

		})
	})
}
