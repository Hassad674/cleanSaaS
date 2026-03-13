export type Referral = {
  id: string;
  referrer_id: string;
  referred_id?: string;
  code: string;
  status: "pending" | "completed" | "rewarded";
  reward_type: "credit" | "trial_extension";
  reward_amount: number;
  created_at: string;
  completed_at?: string;
  rewarded_at?: string;
};

export type ReferralStats = {
  total_referrals: number;
  completed_referrals: number;
  total_rewards: number;
};

export type ReferralCodeResponse = {
  code: string;
};

export type ReferralListResponse = {
  referrals: Referral[];
  total: number;
  page: number;
  limit: number;
};

export type ApplyReferralResponse = {
  message: string;
};
