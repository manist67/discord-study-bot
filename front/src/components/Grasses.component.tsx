import dayjs from "dayjs";
import "./Grasses.component.scss"
import { useState } from "react";

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
  } else if (duration > 0) {
    return 1;
  }
  return 0;
}

export function Grasses({
  data
}: GrassesProps) {
  const [tooltip, setTooltip] = useState<Tooltip>({ show: false, x: 0, y: 0, content: '' });
  const startDate = dayjs().subtract(364, "day");
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
      <div className="grasses-wrapper">
        {tooltip.show && (
          <div className="tooltip" style={{ left: tooltip.x, top: tooltip.y }}>
            {tooltip.content}
          </div>
        )}
        {grasses.map((date) => {
          const duration = data.get(date) ?? 0;
          const level = getLevel(duration);
          return (
            <div
              key={`grass-${date}`}
              className={`grass level-${level}`}
              onMouseOver={(e) => handleMouseOver(e, `${date}: ${Math.floor(duration / 60)} minutes`)}
              onMouseOut={handleMouseOut}
            />
          )
        })}
      </div>
      <div className="grasses-tips-wrapper">
        <span>Less</span>
        <div className="grass level-0"/>
        <div className="grass level-1"/>
        <div className="grass level-2"/>
        <div className="grass level-3"/>
        <div className="grass level-4"/>
        <span>More</span>
      </div>
    </div>
  )
}