import type { MemberResponse } from "../dto/memberResponse.dto";
import { api } from "./api";

export async function getMemeberInfo(guildId: string, memberId: string): Promise<MemberResponse> {
  const { data } = await api.get<MemberResponse>(`${guildId}/${memberId}`)

  return data;
}