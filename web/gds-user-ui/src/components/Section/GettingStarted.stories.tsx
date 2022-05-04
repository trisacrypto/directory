import React from 'react';
import { Story } from '@storybook/react';
import GettingStartedSection from './GettingStarted';

interface GettingStartedProps {}

export default {
  title: 'components/GettingStarted',
  component: GettingStartedSection
};

export const Default: Story<GettingStartedProps> = ({ ...props }) => (
  <GettingStartedSection {...props} />
);

Default.bind({});
