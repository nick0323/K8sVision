
import { FaBars } from 'react-icons/fa';

const PageHeader = ({ title, children, collapsed, onToggleCollapsed }) => {
  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'space-between', 
      alignItems: 'center', 
      marginBottom: 8 
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
        <div className="overview-title-main" style={{ 
          marginBottom: 0, 
          height: 'auto', 
          lineHeight: 'normal' 
        }}>
          {title}
        </div>
        <span style={{ 
          color: '#e0e0e0', 
          fontSize: '18px', 
          fontWeight: '300',
          margin: '0 2px'
        }}>
          |
        </span>
        <button 
          className="page-collapse-btn" 
          onClick={onToggleCollapsed} 
          aria-label="Toggle sidebar"
        >
          <FaBars />
        </button>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        {children}
      </div>
    </div>
  );
};

export default PageHeader; 