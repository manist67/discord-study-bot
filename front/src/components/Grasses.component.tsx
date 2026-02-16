import dayjs from "dayjs";
import "./Grasses.component.scss"
import { useState } from "react";
import { formatDuration } from "../utils/dateString";

type GrassesProps = {
  data: Map<string, number>
}

type Tooltip = {
  show: boolean;
  x: number;
  y: number;
  content: string;
}

const getLevel = (duration: number) => {
  if (duration >= 6400) {
    return 4;
  } else if (duration >= 3600) {
    return 3;
  } else if (duration >= 1800) {
    return 2;
  } else if (duration > 300) {
    return 1;
  }
  return 0;
}

const MONTHES = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec"
]

const WEEKS = [
  "일",
  "월",
  "화",
  "수",
  "목",
  "금",
  "토",
]

export function Grasses({
  data
}: GrassesProps) {
  const [tooltip, setTooltip] = useState<Tooltip>({ show: false, x: 0, y: 0, content: '' });
  const startDate = dayjs().subtract(363, "day");
  const monthItems = Array.from({ length: 52 }, (_, idx) => {
    return MONTHES[startDate.add(idx, "week").month()]
  }).map((e, idx, arr) => {
    if(idx == 0) return e;
    if(e == arr[idx-1]) return null;
    return e
  });
  const weeksItems = Array.from({ length: 7 }, (_,idx) => {
    return WEEKS[startDate.add(idx, "d").day()]
  }).map((e, idx) => {
    if(idx == 0) return e;
    if(idx == 3) return e;
    if(idx == 6) return e;
    return null
  });

  const grasses = Array.from({ length: 364 }, (_, idx) => {
    return startDate.add(idx, "day").format("YYYY-MM-DD");
  });
  const handleMouseOver = (e: React.MouseEvent<HTMLDivElement>, content: string) => {
    const rect = e.currentTarget.getBoundingClientRect();
    setTooltip({
      show: true,
      x: rect.x + window.scrollX + rect.width / 2,
      y: rect.y + window.scrollY - 30,
      content: content,
    });
  };

  const handleMouseOut = () => {
    setTooltip({ ...tooltip, show: false });
  };

  return (
    <div className="Grasses">
      <div className="month-viewer">
        {monthItems.map((e, idx)=>{
          if(e == null) return <span key={`week-${idx}`}/>
          return <span key={`week-${idx}`}>{e}</span>
        })}
      </div>
      <div className="week-wrapper">
        <div className="week-viewer">
        {weeksItems.map(((e, idx) => {
          if(e == null) return <span key={`day-of-week-${idx}`}/>
          return <span key={`day-of-week-${idx}`}>{e}</span>
        }))}
        </div>
        <div className="grasses-wrapper">
          {grasses.map((date) => {
            const duration = data.get(date) ?? 0;
            const level = getLevel(duration);
            return (
              <div
                key={`grass-${date}`}
                className={`grass level-${level}`}
                onMouseOver={(e) => handleMouseOver(e, `${date}: ${formatDuration(duration)}`)}
                onMouseOut={handleMouseOut}
              />
            )
          })}
        </div>
      </div>
      <div className="grasses-tips-wrapper">
        <span>Less</span>
        <div className="grass level-0" 
          onMouseOver={(e) => handleMouseOver(e, "Less then 5m")}
          onMouseOut={handleMouseOut}/>
        <div className="grass level-1"
          onMouseOver={(e) => handleMouseOver(e, "Less then 30m")}
          onMouseOut={handleMouseOut}/>
        <div className="grass level-2"
          onMouseOver={(e) => handleMouseOver(e, "Less then 1h")}
          onMouseOut={handleMouseOut}/>
        <div className="grass level-3"
          onMouseOver={(e) => handleMouseOver(e, "Less then 2h")}
          onMouseOut={handleMouseOut}/>
        <div className="grass level-4"
          onMouseOver={(e) => handleMouseOver(e, "More then 2h")}
          onMouseOut={handleMouseOut}/>
        <span>More</span>
      </div>
      {tooltip.show && (
        <div className="tooltip" style={{ left: tooltip.x, top: tooltip.y }}>
          {tooltip.content}
        </div>
      )}
    </div>
  )
}