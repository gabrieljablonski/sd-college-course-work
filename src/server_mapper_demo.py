import random
import matplotlib.pyplot as plt

def plot_connections(server_number):
    response = {}
    base_delta = int(round(number_of_servers**.5))
    sX = server_number % base_delta
    sY = server_number//base_delta
    # north
    tX, tY = sX, sY
    b = 1
    while 1:
        tY = sY - b
        if tY < 0:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # south
    tX, tY = sX, sY
    b = 1
    while 1:
        tY = sY + b
        if tY >= base_delta:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # west
    tX, tY = sX, sY
    b = 1
    while 1:
        tX = sX - b
        if tX < 0:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # east
    tX, tY = sX, sY
    b = 1
    while 1:
        tX = sX + b
        if tX >= base_delta:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # northeast
    tX, tY = sX, sY
    b = 1
    while 1:
        tX = sX + b
        tY = sY - b
        if tX >= base_delta or tY < 0:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # southeast
    tX, tY = sX, sY
    b = 1
    while 1:
        tX = sX + b
        tY = sY + b
        if tX >= base_delta or tY >= base_delta:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # southwest
    tX, tY = sX, sY
    b = 1
    while 1:
        tX = sX - b
        tY = sY + b
        if tX < 0 or tY >= base_delta:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2
    # northwest
    tX, tY = sX, sY
    b = 1
    while 1:
        tX = sX - b
        tY = sY - b
        if tX < 0 or tY < 0:
            break
        tN = tY*base_delta + tX
        response[str(tN)] = ip_list[tN]
        b *= 2

    print(response)

    fig = plt.figure(frameon=False)
    ax = fig.add_axes([0,0,1,1])
    ax.axis('off')

    xs = list(range(base_delta))
    ys = list(range(base_delta))

    marker_size = 1
    big_marker_size = 10

    for x in xs:
        plt.plot([x]*base_delta, ys, color='black', linestyle='None', marker='o', markersize=marker_size)

    for r in response.values():
        x = r%base_delta
        y = base_delta-1-r//base_delta
        plt.plot([x],[y], 'bo', markersize=big_marker_size)

    plt.plot([sX], [base_delta-1-sY], color='red', linestyle='None', marker='x', markersize=big_marker_size)

    with open(f"_server_connections/conn_{server_number}.png", 'wb') as f:
        fig.canvas.print_png(f)

number_of_servers = 50**2
ip_list = list(range(number_of_servers))

for _ in range(10):
    plot_connections(random.randint(0, number_of_servers-1))
