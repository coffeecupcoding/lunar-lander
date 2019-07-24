function f330 () L=L+S; T=T-S; M=M-S*K; A=I; V=J end
function f420 () Q=S*K/M; J=V+G*S+Z*(-Q-Q*Q/2-Q^3/3-Q^4/4-Q^5/5)
I=A-G*S*S/2-V*S+Z*S*(Q/2+Q^2/6+Q^3/12+Q^4/20+Q^5/30) end
print("               LUNAR")
print("CREATIVE COMPUTING MORRISTOWN, NEW JERSEY")
print(); print(); print()
print( "THIS IS A COMPUTER SIMULATION OF AN APOLLO LUNAR")
print( "LANDING CAPSULE."); print(); print()
print( "THE ON-BOARD COMPUTER HAS FAILED (IT WAS MADE BY")
print( "XEROX) SO YOU HAVE TO LAND THE CAPSULE MANUALLY.")
::line70:: print(); print( "SET BURN RATE OF RETRO ROCKETS TO ANY VALUE BETWEEN")
print( "0 (FREE FALL) AND 200 (MAXIMUM BURN) POUNDS PER SECOND.")
print( "SET NEW BURN RATE EVERY 10 SECONDS."); print()
print( "CAPSULE WEIGHT 32,500 LBS; FUEL WEIGHT 16,500 LBS.")
print(); print(); print(); print( "GOOD LUCK")
L=0
print(); print( "SEC","MI + FT","MPH","LB FUEL","BURN RATE")
A=120; V=1; M=33000; N=16500; G=0.001; Z=1.8
::line150:: io.write( L .. "     " .. math.floor(A) .. "   " .. math.floor(5280*(A%1)) .. "   " .. math.floor(3600*V) .. "   " .. M-N .. "  "); K=io.read(); T=10 
::line160:: if M-N<0.001 then goto line240 end
if T<0.001 then goto line150 end
S=T; if M>=N+S*K then goto line200 end
S=(M-N)/K
::line200:: f420(); if I<=0 then goto line340 end
if V<=0 then goto line230 end
if J<0 then goto line370 end
::line230:: f330(); goto line160
::line240:: print("FUEL OUT AT " .. L .. " SECONDS"); S=(-V+math.sqrt(V*V+2*A*G))/G
V=V+G*S; L=L+S
::line260:: W=3600*V; print( "ON MOON AT " .. L .. " SECONDS - IMPACT VELOCITY " .. W .. " MPH")
if W<=1.2 then print( "PERFECT LANDING!"); goto line440 end
if W<=10 then print( "GOOD LANDING (COULD RE BETTER)"); goto line440 end
if W>60 then goto line300 end
print( "CRAFT DAMAGE... YOU'RE STRANDED HERE UNTIL A RESCUE")
print( "PARTY ARRIVES. HOPE YOU HAVE ENOUGH OXYGEN!")
goto line440
::line300:: print( "SORRY THERE WERE NO SURVIVORS. YOU BLEW IT!")
print( "IN FACT, YOU BLASTED A NEW LUNAR CRATER " .. W*.227 .. "FEET DEEP!")
goto line440
::line340:: if S<0.005 then goto line260 end
D=V+math.sqrt(V*V+2*A*(G-Z*K/M)); S=2*A/D
f420(); f330(); goto line340
::line370:: W=(1-M*G/(Z*K))/2; S=M*V/(Z*K*(W+math.sqrt(W*W+V/Z)))+.05; f420()
if I<=0 then goto line340 end
f330(); if J>0 then goto line160 end
if V>0 then goto line370 end
goto line160
::line440:: print(); print(); print(); print( "TRY AGAIN??"); goto line70
