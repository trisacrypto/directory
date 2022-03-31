import { Meta, Story } from '@storybook/react';
import ReviewSubmit from '.';

export default {
  title: 'components/ReviewSubmit',
  component: ReviewSubmit
} as Meta;

const Template: Story = (args) => <ReviewSubmit {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
