import { useEffect, useState } from "react"
import { useParams } from "react-router"
import { GuildResponse, type GuildResponseType } from "../dto/guildResponse.dto";

export function GuildInfo () { 
  const params = useParams()
  const [ guild, setGuild ] = useState<GuildResponseType>();

  useEffect(()=>{
    if(!params.guildId) return;
    async function getData() {
      if(!params.guildId) return;
      const guildId = params.guildId as string;
      console.log(guildId)
      try {
        const response = await fetch(`/api/${guildId}`);
        if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
        }

        const json = await response.json();
        const res = GuildResponse.parse(json)
        setGuild(res)
      } catch (error) {
        console.error(error);
      }
    }

    getData()
  }, [params.guildId])
  
  return <>
    <h1>{guild?.guild.guildName}</h1>
    <h2>총 사용자 활동 시간 : </h2>
    <div className="winner-wrapper">
      <p>가장 많이 한 사람</p>
      <div className="rank-wrapper">
        <span>
          1
        </span>
        <div className="info-wrapper">
          <p>김민성</p>
          <p>이번 달 : 23시간 17분</p>
        </div>
      </div>
    </div>
    <div className="ranks-wrapper">
      {guild?.members.map((member, idx)=>{
        let remainTime = member.time;
        const hour = Math.floor(remainTime / 3600)
        remainTime %= 3600
        const min = Math.floor(remainTime / 60)
        remainTime %= 60

        return (
          <div className="rank-wrapper">
            <span>{idx + 1}</span>
            <div className="info-wrapper">
              <p>{member.memberName}</p>
              <p>이번 달 : {hour}시간 {min}분 {remainTime}초</p>
            </div>
          </div>
        )
      })}
    </div>
  </>
}