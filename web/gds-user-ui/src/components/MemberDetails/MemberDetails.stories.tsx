import { Meta, Story } from '@storybook/react';
import MemberDetails from '.';

export default {
  title: 'components/MemberDetails',
  component: MemberDetails
} as Meta;

const Template: Story = (args) => <MemberDetails {...args} />;

export const Standard = Template.bind({});
