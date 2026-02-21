import { useState } from 'react';
import './App.css';

export default function App() {
  const [start, setStart] = useState('2026-02-01T00:00');
  const [end, setEnd] = useState('2026-02-20T00:00');
  const [all, setAll] = useState(false);
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchAnalysis = async () => {
    setLoading(true);
    try {
      const startIso = new Date(start).toISOString();
      const endIso = new Date(end).toISOString();
      const res = await fetch(`http://localhost:8080/api/v1/analysis?start=${startIso}&end=${endIso}&all=${all}`);
      setData(await res.json() || []);
    } catch { 
      alert("Error fetching data"); 
    }
    setLoading(false);
  };

  return (
    <div className="container">
      <h1>SMARD-METEO API ⚡</h1>
      <div className="card">
        <div className="controls">
          <input type="datetime-local" value={start} onChange={e => setStart(e.target.value)} />
          <input type="datetime-local" value={end} onChange={e => setEnd(e.target.value)} />
          <label><input type="checkbox" checked={all} onChange={e => setAll(e.target.checked)} /> All Data</label>
          <button onClick={fetchAnalysis} disabled={loading}>{loading ? '...' : 'Analyze'}</button>
        </div>
        
        <table>
          <thead>
            <tr>
              <th>Time</th>
              <th>Consumption</th>
              <th>Renewables</th>
              <th>Weather</th>
              <th>Covered</th>
            </tr>
          </thead>
          <tbody>
            {data.map((r, i) => (
              <tr key={i}>
                <td>{new Date(r.grid_state.timestamp).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}</td>
                <td>{Math.round(r.grid_state.consumption)} MWh</td>
                <td>{Math.round(r.grid_state.generation_renewable)} MWh</td>
                <td className="weather">
                  {r.weather_state ? `💨 ${r.weather_state.wind_speed} km/h | ☀️ ${r.weather_state.solar_radiation} W/m²` : 'N/A'}
                </td>
                <td style={{color: r.grid_state.is_covered ? '#34d399' : '#fca5a5', fontWeight: 'bold'}}>
                  {r.grid_state.is_covered ? 'Yes' : 'No'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
