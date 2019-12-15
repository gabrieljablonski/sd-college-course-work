powershell -noexit D:\GitReps\sd-college-course-work\etcd\etcd.exe ^
    --name %1 ^
    --listen-peer-urls %2 ^
    --initial-advertise-peer-urls %2 ^
    --initial-cluster-token %3 ^
    --initial-cluster %4 ^
    --initial-cluster-state new ^
    --listen-client-urls %5 ^
    --advertise-client-urls %5