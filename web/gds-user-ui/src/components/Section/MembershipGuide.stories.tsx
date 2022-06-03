import { Story } from '@storybook/react';
import MembershipGuide from './MembershipGuide';

export default {
  title: 'Components/MembershipGuide',
  component: MembershipGuide
};

export const Standard: Story = (props) => <MembershipGuide {...props} />;

Standard.bind({});
