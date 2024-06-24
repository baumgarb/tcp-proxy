# Usage

Just run the binary without arguments, should be self-explanatory.

```bash
user@host:~ $ ./tcp-proxy
Usage: ./tcp-proxy <src addr> <dest addr> [-v]

Examples: 
       ./tcp-proxy :3000 :5000                         # forwards all incoming TCP connections on port 3000 to 5000.
       ./tcp-proxy :3000 :5000 &                       # forwards all incoming TCP connections on port 3000 to 5000 
                                                       # while running silently in the background.
       ./tcp-proxy :5000 :31774 -v                     # forwards all incoming TCP connections on port 5000 to 31774.
                                                       # Verbose logging is enabled.
       ./tcp-proxy 127.0.0.1:3000 192.168.100.10:5000  # forwards all incoming TCP connections on 127.0.0.1 on port 3000
                                                       # to target with IP 192.168.100.10 on port 5000.
```

# Why on earth ...?

I'm running a K3s in a WSL 2 instance on my local Windows machine. Creating a node port service on K3s is not opening the port on the WSL instance, instead it's using some iptables/ipvs/ebpf/whatever magic and therefore the underlying WSL VM does not know about this port being relevant on the WSL instance. That's expected K8s/K3s behavior. Because of that any web service or web API I'm exposing via a node port service cannot be reached on the Windows host via `http://localhost`. Using `kubectl port-forward` is not an option for my scenario since it opens the port to one specific pod of an underlying deployment, it's not balancing the request across all replicas, but that's what I want/need. An easy way to make it work is to create a port forward on Windows (open console as admin and run `netsh interface portproxy add v4tov4 listenaddress=0.0.0.0 listenport=<src port> connectaddress=<WSL VM IP> connectport=<dest port>`), but that didn't work consistenly in my case because I'm running multiple WSL instances based on different Linux distributions. And all of them are "hiding" behind one IP address you can use for the `netsh` rule. Simply run `ip addr show eth0` on each instance and you'll see that all of them have the same IP address. Once you have a couple of WSL instances open it stops working when you want to expose new services. 

To make it work I was searching for the most basic and simple TCP proxy that essentially lets the WSL VM know about an open port on any of my instances. And it should come with as few as possible security and vulnerability attach vectors. On top of that, I didn't want to setup a fancy production-grade proxy solution (Nginx, Apache, Traefik, ...) which comes with all the bells and whistles I just don't need for my purposes. 

Therefore, I created this super simple, dumb and inefficient TCP proxy. I know what's in there, no additional packages or libraries pulled in, it's all Go's standard library. And for my simple purposes it is okay that it does not respect something like a `Connection: keep-alive` HTTP header. It's fine for me to open a separate TCP connection for each and every HTTP resource requested. After all, it's just a local dev playground thingy and it solves my problem.