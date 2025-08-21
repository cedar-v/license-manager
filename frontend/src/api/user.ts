import Axios from './https'

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  code: string;
  message: string;
  data: {
    token?: string;
    [key: string]: any;
  };
}

export function Login(data: LoginRequest): Promise<LoginResponse> {
  return Axios.post('/api/v1/login', data)
}