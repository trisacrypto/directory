const MENU_ITEMS = [
    { key: 'registration', label: 'Review Registration', isTitle: true },
    { key: 'dashboard', icon: 'uil-home-alt', label: 'Dashboard', url: '/dashboard' },

    { key: 'vasps', label: 'VASPs Summary', isTitle: true },
    { key: 'vs-list', icon: 'uil-list-ul', label: 'List', url: '/vasps-summary/vasps' },
    // { key: 'vs-details', icon: 'uil-folder-question', label: 'Details', url: '/vasps-summary/vasps/:id/details' },
];

export default MENU_ITEMS;
