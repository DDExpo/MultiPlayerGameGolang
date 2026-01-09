type WSMessage<T = any> = {
    type: string
    data: T
}