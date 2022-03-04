import React from 'react';
import { Story } from '@storybook/react';
import VaspVerification from './VaspVerification';

interface VaspVerificationProps {}

export default {
  title: 'Components/VaspVerification',
  component: VaspVerification
};

export const standard: Story<VaspVerificationProps> = ({ ...props }) => (
  <VaspVerification {...props} />
);

standard.bind({});
