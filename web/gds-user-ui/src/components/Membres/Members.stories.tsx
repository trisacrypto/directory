import { Meta, Story } from '@storybook/react';

import Members from '.';

type MembersProps = {
  key?: string;
  name: string;
  isTestNet: boolean;
  isMainNet: boolean;
  memberId: string;
};

export default {
  title: 'components/Members',
  component: Members
} as Meta;

const Template: Story<MembersProps> = (args) => <Members {...args} />;

export const Default = Template.bind({});
Default.args = {};
