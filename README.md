# Nextime

This small tool will read a nextpnr report file and show information about
timing and utilization. By default it will just show utilization. If the report
indicates that timing has failed, Nextime will also display the critical path
that caused it to fail.

# Example

```
$ nextpnr-ecp5 --report report.json --freq 48 ...
...
ERROR: Max frequency for clock   '$glbnet$clk_sys': 44.71 MHz (FAIL at 48.00 MHz)
$ nextime -util report.json
CELL         	TOTAL	USED	UTILIZATION 
DCCA         	56   	1   	0.02       	
DP16KD       	56   	24  	0.43       	
EHXPLLL      	2    	1   	0.50       	
TRELLIS_IO   	197  	6   	0.03       	
TRELLIS_SLICE	12144	1890	0.16       	

TOTAL	ROUTING	LOGIC	NET                                                    
0.00 	0.52   	0.00 	top.core.dpath.rf.regs.0.1.0$D                        	
0.52 	1.15   	0.24 	top.core.dpath.rf.regs.0.0.0_RAD[3]                   	
1.91 	0.89   	0.24 	top.core.dpath.rf.regs.0.0.0_DO_3[1]                  	
3.04 	0.73   	0.25 	top.core.dpath.brCond.io_rs2_L6MUX21_Z_SD[6]          	
4.02 	1.54   	0.26 	top.core.ctrl._signals_T_114_LUT4_D_Z_PFUMX_Z_2_C0[1] 	
5.81 	0.00   	0.24 	top.core.dpath.alu._io_out_T_4_LUT4_Z_7_D_L6MUX21_Z_D1	
6.06 	1.00   	0.24 	top.core.dpath.alu._io_out_T_4_LUT4_Z_7_D[6]          	
7.29 	1.17   	0.45 	top.core.dpath.alu_io_b[0]                            	
8.91 	0.00   	0.07 	top.core.dpath.alu._b_T_1_CCU2C_S0_3_COUT[1]          	
8.98 	0.00   	0.07 	top.core.dpath.alu._b_T_1_CCU2C_S0_3_COUT[3]          	
9.05 	0.00   	0.47 	top.core.dpath.alu._b_T_1_CCU2C_S0_3_COUT[5]          	
9.52 	0.90   	0.24 	top.core.dpath.alu._b_T_1[7]                          	
10.66	1.16   	0.45 	top.core.dpath.alu.b[7]                               	
12.27	0.00   	0.07 	top.core.dpath.alu.io_sum_CCU2C_S0_3_COUT[7]          	
12.34	0.00   	0.07 	top.core.dpath.alu.io_sum_CCU2C_S0_3_COUT[9]          	
12.41	0.00   	0.07 	top.core.dpath.alu.io_sum_CCU2C_S0_3_COUT[11]         	
12.48	0.00   	0.07 	top.core.dpath.alu.io_sum_CCU2C_S0_3_COUT[13]         	
12.55	0.00   	0.07 	top.core.dpath.alu.io_sum_CCU2C_S0_3_COUT[15]         	
12.62	0.00   	0.47 	top.core.dpath.alu.io_sum_CCU2C_S0_3_COUT[17]         	
13.10	0.88   	0.24 	top.bus_io_host_0_addr[19]                            	
14.22	0.67   	0.26 	top.bus.devSelReq_LUT4_Z_A_PFUMX_Z_C0[4]              	
15.14	1.75   	0.24 	top.bus.devSelReq_LUT4_Z_A[1]                         	
17.12	0.64   	0.24 	top.bus.devSelReq                                     	
18.00	0.69   	0.26 	top.bus_io_dev_0_req                                  	
18.94	1.00   	0.24 	top.ram._write_T_17_LUT4_Z_1_B[2]                     	
20.17	1.97   	0.22 	top.ram._GEN_15[23]                                   	
22.15	0.22   	0.00 	top.ram.mem.5.0.0[DIA3]                               	
Critical path: top.core.dpath.rf.regs.0.1.0$D -> top.ram.mem.5.0.0[DIA3]
Max frequency: 44.71 MHz (22.36 ns)
$glbnet$clk_sys failed at 48.00 MHz
```

# Usage

```
Usage of nextime:
  -crit string
    	show critical path for clock
  -util
    	show utilization
```

If no JSON report is provided, Nextime automatically tries to use `report.json`.

When Nextime runs with no arguments it will only report a critical path if
timing is not met. To view a critical path when timing is met you must request
it with `-crit`.
