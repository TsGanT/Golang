import sys
import random
sys.setrecursionlimit(10000)

def big_odd(b):
    num1 = pow(2,1023)+1
    res = num1+b*2
    return res

def Fermat(x, n, p):
    if n==0:
        return 1
    res = Fermat((x*x)%p, n>>1, p)
    if n&1 != 0:
        res = (res*(x))%p
    return res

def MillerRabin(a,p):
    if Fermat(a,p-1,p)==1:
        u = (p-1) >> 1
        while u&1 == 0:
            t = Fermat(a,u,p)
            if t == 1:
                u = u >> 1
            else:
                if t == p-1:
                    return True
                else:
                    return False
        else:
            t = Fermat(a,u,p)
            if t == 1 or t == p-1:
                return True
            else:
                return False
    else:       #This is not very useful, it may cause more output
        return False

def TestMillarRabin(p, root):
    for k in range(0,6):
        if not MillerRabin(root,p):        # Here is the place to set your root!
            return False
    print("Success!! You have found a prime number/0.0002")
    return p 

def GetPrime():
    seed = random.randint(0,5000)
    for i in range(seed, seed+1000):
        resp = big_odd(i)
        res2 = TestMillarRabin(resp, 7)
        if res2:
            #if checkbit(res2):
                return res2

def checkbit(p):
    count = 0
    while p>0:
        p = p>>1
        count = count + 1
    if count == 1024:
        return True
    else:
        return False
p = GetPrime()
print("This is the 1024 prime:",p)