import { Meta, Story } from '@storybook/react';
import BasicDetails from './';

type BasicDetailsProps = {};

export default {
  title: 'components/BasicDetails',
  component: BasicDetails
} as Meta<BasicDetailsProps>;

const Template: Story<BasicDetailsProps> = (args) => <BasicDetails {...args} />;

export const Default = Template.bind({});
Default.args = {};
