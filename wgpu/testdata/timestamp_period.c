#if defined(_WIN32)
#define EXPORT __declspec(dllexport)
#else
#define EXPORT __attribute__((visibility("default")))
#endif

EXPORT float wgpuQueueGetTimestampPeriod(const void *queue) {
    return queue ? 0.125f : 0.0f;
}
