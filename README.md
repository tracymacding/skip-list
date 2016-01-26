# skip-list

## Performance

### Single-Thread Write

~ 437896 TPS

### Single-Thread Read

Use Lock: ~ 676997 TPS
Non Lock: ~

### Multi-Thread Write/Read 

Read:Write 1:1  
Non Lock:  
Read:  ~ 621703 TPS
Write: ~ 498162 TPS

Use Lock:  
Read:  ~ 493449 TPS
Write: ~ 221596 TPS

Read:Write 5:10
Non Lock:  
Read:  ~ 435806 TPS
Write: ~ 41114  TPS

Use Lock:  
Read:  ~ 69749 TPS
Write: ~ 33660 TPS

Read:Write 10:5
Non Lock:  
Read:  ~ 356519 TPS
Write: ~ 101703 TPS

Use Lock:  
Read:  ~ 46776  TPS
Write: ~ 78991  TPS
