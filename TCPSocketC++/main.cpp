#include "socket.hpp"

int main()
{
    Socket socket = Socket(80, INADDR_ANY);
    if (socket.success == -1)
    {
        return EXIT_FAILURE;
    }
    while (true)
    {
        socket.startListening();
        if (socket.success == -1)
        {
            return EXIT_FAILURE;
        }
        std::cout << "Response done\n";
    }

    return 0;
}