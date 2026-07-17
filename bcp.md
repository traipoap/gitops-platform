Phase 1: Pre-Migration Verification

brk-edge-02 (Old Gateway)
show running-config interface port-channel 1.535
show ip interface brief | include port-channel 1.535
show ip arp interface port-channel 1.535
show ip route bgp | include 164.115.47.0
ping 8.8.8.8 source port-channel 1.535
show ip arp

nbi-agg-01 (New Gateway)
show running-config | include router bgp
show ip interface brief
show ip route bgp

Phase 2: Safe Staging (เตรียมคอนฟิกแบบปลอดภัย)

nbi-agg-01 (New Gateway)
configure terminal
interface port-channel 1.535
 description #### Workload ####
 encapsulation dot1q 535
 ip address 164.115.47.1 255.255.255.0
 shutdown
 exit
 
Phase 3: The Cutover Sequence (ขั้นตอนตัดสลับ)

ปิดการทำงานของ Gateway เดิมบน brk-edge-02
configure terminal
interface port-channel 1.535
 shutdown
 exit
 
เปิดใช้งาน Gateway ใหม่บน nbi-agg-01
configure terminal
interface port-channel 1.535
 no shutdown
 exit

ถอนการประกาศ BGP บน brk-edge-02
configure terminal
router bgp 135566
 address-family ipv4
  no network 164.115.47.0 mask 255.255.255.0
  exit
 exit
 
เริ่มประกาศ BGP บน nbi-agg-01
configure terminal
router bgp 9835
 address-family ipv4
  network 164.115.47.0 mask 255.255.255.0
  exit
 exit

Phase 4: Post-Migration Verification (ตรวจสอบหลังย้าย)
ตรวจสอบบน nbi-agg-01 (New Gateway)
ping 8.8.8.8 source 164.115.47.1
show ip arp interface port-channel 1.535
show ip route bgp | include 164.115.47.0

ตรวจสอบความเรียบร้อยก่อนบันทึกค่า (Save Config)
บน brk-edge-02 (ลบ Sub-interface และเซฟ):
configure terminal
no interface port-channel 1.535
end
write memory

บน nbi-agg-01 (เซฟคอนฟิก):
end
write memory

Phase 5: Complete Rollback Plan (แผนกู้คืนระบบแบบร้อยเปอร์เซ็นต์)
Step 1: ปิดการทำงานบน nbi-agg-01 (New Gateway)
configure terminal
interface port-channel 1.535
 shutdown
 exit
 
router bgp 9835
 address-family ipv4
  no network 164.115.47.0 mask 255.255.255.0
  exit
 exit
end

Step 2: เปิดการทำงานบน brk-edge-02 (Old Gateway)
configure terminal
interface port-channel 1.535
 no shutdown
 exit
 
router bgp 135566
 address-family ipv4
  network 164.115.47.0 mask 255.255.255.0
  exit
 exit
end

Step 3: ตรวจสอบและดึงทราฟฟิกกลับมา
ping 8.8.8.8 source port-channel 1.535
show ip arp interface port-channel 1.535
