export type Plan = {
  id: string;
  name: string;
  price: number;
  currency: string;
  interval: "month" | "year";
  features: string[];
};

export type Subscription = {
  id: string;
  plan_id: string;
  status: "active" | "canceled" | "past_due" | "trialing" | "inactive";
  current_period_end: string;
  cancel_at_period_end: boolean;
};

export type Invoice = {
  id: string;
  amount: number;
  currency: string;
  status: string;
  created_at: string;
};
