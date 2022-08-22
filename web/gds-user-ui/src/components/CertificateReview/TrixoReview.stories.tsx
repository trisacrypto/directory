import { Meta, Story } from '@storybook/react';

import TrixoReview from './TrixoReview';
interface TrixoReviewProps {
  data: any;
}
export default {
  title: 'components/TrixoReview',
  component: TrixoReview
} as Meta;

const Template: Story<TrixoReviewProps> = (args) => <TrixoReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
