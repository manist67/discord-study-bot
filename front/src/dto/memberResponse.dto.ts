import { z } from 'zod';

export const Particiapting = z.object({
  date: z.iso.datetime(),
  duration: z.number()
})

export const MemberResponse = z.object({
  member: z.object({
    nickname: z.string()
  }),
  total: z.number(),
  weekTotal: z.number(),
  participatingList: z.array(Particiapting)
})

export type Particiapting = z.infer<typeof Particiapting>
export type MemberResponse = z.infer<typeof MemberResponse>