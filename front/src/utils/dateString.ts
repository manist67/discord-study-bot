export const formatDuration = (seconds: number) => {
  if(seconds == 0) return "0s"
  const text: string[] = [];
  if (seconds >= 3600) {
    text.push(`${Math.floor(seconds / 3600)}h`);
    seconds = seconds % 3600;
  }

  if(seconds >= 60) {
    const remainMin = Math.floor(seconds / 60)
    if(remainMin != 0) {
      text.push(`${remainMin}m`)
    }

    seconds = seconds % 60
  }

  if(seconds > 0) {
    text.push(`${seconds}s`)
  }
  
  return text.join(" ")
};