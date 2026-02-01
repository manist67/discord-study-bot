import { Route, Routes } from 'react-router'
import './App.css'
import { Home } from './pages/Home'
import { GuildInfo } from './pages/GuildInfo'
import { MemberInfo } from './pages/MemberInfo'

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/:guildId" element={<GuildInfo />} />
      <Route path="/:guildId/:memberId" element={<MemberInfo />} />
    </Routes>
  )
}

export default App
