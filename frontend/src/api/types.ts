export interface ApiSuccessEnvelope<T> {
  timestamp: string
  code: number
  message: string
  data?: T
}

export interface ApiErrorEnvelope {
  timestamp: string
  code: number
  message: string
  error?: unknown
}
