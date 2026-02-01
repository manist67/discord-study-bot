import { z } from 'zod';

export const GuildResponse = z.object({
  guild: z.object({
    idx: z.int(),
    guildId: z.string(),
    guildName: z.string(),
  }),
  members: z.array(z.object({
    time: z.int(),
    memberId: z.string(),
    memberName: z.string(),
  }))
})

export type GuildResponseType = z.infer<typeof GuildResponse>