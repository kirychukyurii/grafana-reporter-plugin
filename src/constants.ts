import pluginJson from './plugin.json';
import { NavModelItem } from '@grafana/data';

export const PLUGIN_ID = `${pluginJson.id}`;
export const PLUGIN_BASE_URL = `/a/${PLUGIN_ID}`;

export enum ROUTES {
  SCHEDULES = 'schedules',
  CONFIG = 'config'
}

export const NAVIGATION_TITLE = 'Excel report e-mail scheduler';
export const NAVIGATION_SUBTITLE = `Generate Excel reports from mSupply dashboard. Send the reports to custom created user-groups on pre-defined schedule.`;

// Add a navigation item for each route you would like to display in the navigation bar
export const NAVIGATION: Record<string, NavModelItem> = {
  [ROUTES.SCHEDULES]: {
    id: ROUTES.SCHEDULES,
    text: 'Schedules',
    icon: 'times',
    url: `${PLUGIN_BASE_URL}/schedules`,
  },
  [ROUTES.CONFIG]: {
    id: ROUTES.CONFIG,
    text: 'Configuration',
    icon: 'cog',
    url: `plugins/${PLUGIN_ID}`
  }
};
