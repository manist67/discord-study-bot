import axios from "axios";
import type { MemberResponse } from "../dto/memberResponse.dto";

const api = axios.create({
  baseURL: "http://localhost:8080/api"
})

export async function getMemeberInfo(guildId: string, memberId: string): Promise<MemberResponse> {
  const { data } = await api.get<MemberResponse>(`${guildId}/${memberId}`)

  return data;
}