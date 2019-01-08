//Algorithm.h
#ifndef ALGORITHM_DLL_H
#define ALGORITHM_DLL_H

namespace Algorithm {
    extern "C" __declspec(dllexport) int add(int a, int b);
    extern "C" __declspec(dllexport) int sub(int a, int b);
}

#endif // !ALGORITHM_DLL_H
