// static/js/main.js
console.log("网站加载完成！");

function updateClock() {
    const now = new Date();
    
    // 更新时间
    const time = now.toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
    document.getElementById('clock').textContent = time;
    
    // 更新日期
    const date = now.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'long'
    });
    document.getElementById('date').textContent = date;
}

// 页面加载完成后启动时钟
document.addEventListener('DOMContentLoaded', function() {
    // 立即更新一次
    updateClock();
    // 每秒更新一次
    setInterval(updateClock, 1000);
});
