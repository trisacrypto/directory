import { Meta, Story } from '@storybook/react';
import BasicDetailsReview from './BasicDetailsReview';

interface BasicDetailsReviewProps {
  data: any;
}
export default {
  title: 'components/BasicDetailsReview',
  component: BasicDetailsReview
} as Meta;

const Template: Story<BasicDetailsReviewProps> = (args) => <BasicDetailsReview {...args} />;

export const Default = Template.bind({});
Default.args = {};
