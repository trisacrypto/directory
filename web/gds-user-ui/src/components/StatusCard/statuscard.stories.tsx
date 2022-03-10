import { Meta, Story } from "@storybook/react";
import StatusCard from ".";

export default {
  title: "components/StatusCard",
  component: StatusCard,
} as Meta;

const Template: Story = (args) => <StatusCard {...args} />;

export const Default = Template.bind({});
Default.args = {};
