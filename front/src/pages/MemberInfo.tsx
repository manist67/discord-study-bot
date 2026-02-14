import { useEffect, useMemo, useState } from "react"
import { getMemeberInfo } from "../api/member"
import { useParams } from "react-router"
import type { MemberResponse } from "../dto/memberResponse.dto";
import dayjs from "dayjs";
import { Grasses } from "../components/Grasses.component";
import "./MemberInfo.scss";
import { formatDuration } from "../utils/dateString";

const INITIAL_VISIBLE_COUNT = 30;

export function MemberInfo() {
  const params = useParams();
  const [response, setResponse] = useState<MemberResponse | null>();
  const [visibleCount, setVisibleCount] = useState<number>(INITIAL_VISIBLE_COUNT);

  useEffect(() => {
    if (!params.guildId || !params.memberId) return;

    getMemeberInfo(params.guildId, params.memberId).then(e => {
      setResponse(e)
    }).catch(e => {
      console.log(e)
    })
  }, [params.guildId, params.memberId]);

  const grassData = useMemo(() => {
    const map = new Map<string, number>()
    if (!response) return map;
    response.participatingList.forEach(v => {
      map.set(dayjs(v.date).format("YYYY-MM-DD"), v.duration);
    })
    return map;
  }, [response]);

  const sortedParticipatingList = useMemo(() => {
    if (!response) return [];
    return [...response.participatingList].sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
  }, [response]);

  const handleLoadMore = () => {
    setVisibleCount(prevCount => prevCount + 50);
  };

  return (
    <div className="member-info-page">
      {response && (
        <>
          <div className="profile-header">
            <div className="avatar-placeholder">
              {response.member.nickname.charAt(0).toUpperCase()}
            </div>
            <div className="info-wrapper">
              <div className="name-wrapper">
                <h1>{response.member.nickname}</h1>
                <span className={`status-label ${response.isOnline ? 'status-online' : 'status-offline'}`}>
                  <span className="status-dot"></span>
                  {response.isOnline ? 'Online' : 'Offline'}
                </span>
              </div>
              <p className="member-name">{response.member.memberName}</p>
            </div>
          </div>

          <div className="stats-section">
            <div className="stat-card">
              <span className="stat-label">총 시간</span>
              <span className="stat-value">{formatDuration(response.total)}</span>
            </div>
            <div className="stat-card">
              <span className="stat-label">이번 주 시간</span>
              <span className="stat-value">{formatDuration(response.weekTotal)}</span>
            </div>
            <div className="stat-card">
              <span className="stat-label">총 활동일</span>
              <span className="stat-value">{response.participatingList.length}일</span>
            </div>
          </div>

          <div className="activity-section">
            <h2>활동 내역</h2>
            <div className="grasses-container">
              <Grasses data={grassData} />
            </div>
          </div>

          <div className="detailed-activity-section">
            <h2>일별 접속 기록</h2>
            <div className="activity-list">
              {sortedParticipatingList.slice(0, visibleCount).map((e, index) => (
                <div className="activity-item" key={index}>
                  <span className="date">{dayjs(e.date).format("YYYY-MM-DD")}</span>
                  <span className="duration">{formatDuration(e.duration)}</span>
                </div>
              ))}
            </div>
            {visibleCount < sortedParticipatingList.length && (
              <div className="load-more-wrapper">
                <button onClick={handleLoadMore}>더 보기</button>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}