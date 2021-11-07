#include <iostream>
#include <unistd.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <netinet/in.h>
#include <string>
#include <vector>

class Socket
{
private:
    int s;
    struct HTTPRequest
    {
        HTTPRequest(std::string input)
        {
            parseString(input);
        }

        std::string verb;
        std::string endpoint;
        std::string host;

    private:
        void parseString(std::string input)
        {
        }
    };

public:
    int success = 0;
    Socket(unsigned short port, u_int32_t address)
    {
        this->s = socket(AF_INET, SOCK_STREAM, 0);
        if (s < 0)
        {
            std::cout << "Could not open SOCKET\n";
            success = -1;
            return;
        }
        sockaddr_in server;
        server.sin_family = AF_INET;
        server.sin_port = htons(port);
        server.sin_addr.s_addr = address;

        int result = bind(s, (sockaddr *)&server, sizeof(server));
        if (result < 0)
        {
            std::cout << "Could not BIND\n";
            success = -1;
            return;
        }
    }

    ~Socket()
    {
        close(s);
    }

    void startListening()
    {
        int result = listen(s, 1);
        if (result != 0)
        {
            std::cout << "Failed to listen\n";
            success = -1;
            return;
        }
        sockaddr_in client;
        socklen_t namelen = sizeof(client);
        int ns = accept(s, (struct sockaddr *)&client, &namelen);
        if (ns == -1)
        {
            std::cout << "Failed to accept\n";
            success = -1;
            return;
        }
        char buf[1024];
        result = recv(ns, buf, sizeof(buf), 0);
        if (result == -1)
        {
            std::cout << "Det sket sig att läsa\n";
            success = -1;
            return;
        }

        std::cout << "bytes read: " << result << "\n";

        std::string str;
        str.assign(buf, result);

        std::cout << str << '\n';

        std::string respString = "HTTP/1.1 200 OK\r\n";

        send(ns, respString.data(), respString.size(), 0);
        close(ns);
    }
};