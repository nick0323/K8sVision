import React, { useEffect, useState, useRef } from 'react';
import InfoCard from './InfoCard';
import ResourceSummary from './ResourceSummary';
import { formatDateTime, formatRelativeTime, useFetch } from './utils';
import { FaChartPie, FaDesktop, FaCubes, FaProjectDiagram } from 'react-icons/fa';
import { LuSquareDashed } from 'react-icons/lu';
import { LOADING_TEXT, ERROR_TEXT, EMPTY_TEXT } from './constants';

export default function OverviewPage() {
  // const [data, setData] = useState({});
  // const [loading, setLoading] = useState(false);
  const { data, loading, error } = useFetch('/api/overview');
  const safeData = data || {};
  const leftColRef = useRef(null);
  const [leftColHeight, setLeftColHeight] = useState(undefined);
  useEffect(() => {
    if (leftColRef.current) {
      setLeftColHeight(leftColRef.current.offsetHeight);
    }
  }, [data, loading, error]);

  return (
    <div>
      {loading && <div style={{textAlign:'center',color:'#888',padding:'32px 0'}}>{LOADING_TEXT}</div>}
      {error && <div style={{textAlign:'center',color:'red',padding:'32px 0'}}>{ERROR_TEXT}{error}</div>}
      {!loading && !error && (
        <>
          <div className="overview-title-main">Overview</div>
          <div className="overview-grid">
            <InfoCard
              icon={<FaDesktop />}
              title="Nodes"
              value={safeData.nodeCount || 0}
              status={safeData.nodeCount === 0 ? <div className="center-empty"><span style={{color:'#c0c4cc',fontSize:'var(--font-size-sm)'}}>{EMPTY_TEXT}</span></div> : (
                <span className={safeData.nodeReady === safeData.nodeCount ? 'status-ready' : 'status-failed'}>
                  {safeData.nodeReady === safeData.nodeCount ? 'All Ready' : `${safeData.nodeCount - safeData.nodeReady} Not Ready`}
                </span>
              )}
            />
            <InfoCard
              icon={<FaCubes />}
              title="Pods"
              value={safeData.podCount || 0}
              status={safeData.podCount === 0 ? <div className="center-empty"><span style={{color:'#c0c4cc',fontSize:'var(--font-size-sm)'}}>{EMPTY_TEXT}</span></div> : (
                <span className={safeData.podNotReady === 0 ? 'status-ready' : 'status-failed'}>
                  {safeData.podNotReady === 0 ? 'All Ready' : `${safeData.podNotReady} Not Ready`}
                </span>
              )}
            />
            <InfoCard
              icon={<LuSquareDashed />}
              title="Namespaces"
              value={safeData.namespaceCount || 0}
              status={safeData.namespaceCount === 0 ? <div className="center-empty"><span style={{color:'#c0c4cc',fontSize:'var(--font-size-sm)'}}>{EMPTY_TEXT}</span></div> : (
                <span className="status-ready">All Ready</span>
              )}
            />
            <InfoCard
              icon={<FaProjectDiagram />}
              title="Services"
              value={safeData.serviceCount || 0}
              status={safeData.serviceCount === 0 ? <div className="center-empty"><span style={{color:'#c0c4cc',fontSize:'var(--font-size-sm)'}}>{EMPTY_TEXT}</span></div> : (
                <span className="status-ready">All Ready</span>
              )}
            />
          </div>
          <div className="overview-row2">
            <div
              className="overview-left-col"
              ref={leftColRef}
              style={leftColHeight ? { minHeight: leftColHeight, display: 'flex', flexDirection: 'column', gap: '16px' } : {}}
            >
              <ResourceSummary style={{ flex: 1 }}
                title="CPU Resources"
                requestsValue={safeData.cpuRequests?.toFixed(1) || 0}
                requestsPercent={((safeData.cpuRequests/safeData.cpuCapacity)*100 || 0).toFixed(1)}
                limitsValue={safeData.cpuLimits?.toFixed(1) || 0}
                limitsPercent={((safeData.cpuLimits/safeData.cpuCapacity)*100 || 0).toFixed(1)}
                totalValue={safeData.cpuCapacity?.toFixed(1) || 0}
                availableValue={(safeData.cpuCapacity - safeData.cpuRequests)?.toFixed(1) || 0}
                unit="cores"
              />
              <ResourceSummary style={{ flex: 1 }}
                title="Memory Resources"
                requestsValue={safeData.memoryRequests?.toFixed(1) || 0}
                requestsPercent={((safeData.memoryRequests/safeData.memoryCapacity)*100 || 0).toFixed(1)}
                limitsValue={safeData.memoryLimits?.toFixed(1) || 0}
                limitsPercent={((safeData.memoryLimits/safeData.memoryCapacity)*100 || 0).toFixed(1)}
                totalValue={safeData.memoryCapacity?.toFixed(1) || 0}
                availableValue={(safeData.memoryCapacity - safeData.memoryRequests)?.toFixed(1) || 0}
                unit="GiB"
              />
            </div>
            <div className="overview-event-col">
              <div
                className="overview-event-card resource-summary-card"
                style={{
                  minHeight: leftColHeight,
                  overflowY: 'auto',
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'center',
                }}
              >
                <div className="resource-summary-title">Recent Events</div>
                {safeData.events && safeData.events.length > 0 ? (
                  safeData.events
                    .slice()
                    .sort((a, b) => new Date(b.lastSeen) - new Date(a.lastSeen))
                    .slice(0, 5)
                    .map((e, i) => (
                      <div key={i} style={{display:'flex',alignItems:'flex-start',padding:'0 0 18px 0',borderBottom: i!==4?'1px solid #f0f0f0':'none',marginBottom: i!==4?12:0}}>
                        <div style={{width:28,display:'flex',justifyContent:'center',alignItems:'flex-start',marginTop:2}}>
                          <span style={{display:'inline-block',width:16,height:16,borderRadius:'50%',border:'2px solid #a5b4fc',background:'#fff',marginTop:2}}></span>
                        </div>
                        <div style={{flex:1}}>
                          <div style={{display:'flex',alignItems:'center',marginBottom:2}}>
                            <span className={e.type === 'Warning' ? 'event-type-warning' : 'event-type-normal'} style={{background: e.type === 'Warning' ? '#ffeaea' : '#e6f7ff', color: e.type === 'Warning' ? '#ff4d4f' : '#1890ff',borderRadius:'10px',padding:'2px 10px',fontSize:'var(--font-size-sm)',fontWeight:600,marginRight:8}}
                            >{e.type}</span>
                            <span style={{fontWeight:600,fontSize:'var(--font-size-sm)',color:'#222',marginRight:8}}>{e.reason}</span>
                            <span style={{marginLeft:'auto',fontSize:'var(--font-size-sm)',color:'#888',fontWeight:400}}>{formatRelativeTime(e.lastSeen)}</span>
                          </div>
                          <div style={{fontSize:'var(--font-size-sm)',color:'#444',marginBottom:2,wordBreak:'break-all'}}>{e.message}</div>
                          <div style={{fontSize:'var(--font-size-sm)',color:'#888',fontWeight:400}}>
                            {e.pod ? `Pod: ${e.pod}` : ''}
                            {e.cloneset ? `CloneSet: ${e.cloneset}` : ''}
                            {e.namespace && e.name ? `Pod: ${e.namespace}/${e.name}` : ''}
                            {e.reporter ? ` Reporter: ${e.reporter}` : ''}
                          </div>
                        </div>
                      </div>
                    ))
                ) : (
                  <div style={{color:'#888',fontSize:'var(--font-size-sm)',padding:'24px 0',textAlign:'center'}}>{EMPTY_TEXT}</div>
                )}
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
} 