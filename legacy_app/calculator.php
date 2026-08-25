<?php
// Legacy Code - No Tests, Messy Logic

class TaxCalculator {
    public function calculateFinalPrice($price, $userType, $discountCode) {
        if ($price < 0) return 0;
        
        $discount = 0;
        if ($discountCode == 'SAVE20') {
            $discount = 0.20;
        } elseif ($discountCode == 'HALF') {
            $discount = 0.50;
        }

        $tax = 0.18; // 18% standard tax
        if ($userType == 'B2B') {
            $tax = 0.05; // 5% for B2B
        } elseif ($userType == 'VIP') {
            $tax = 0.0; // No tax for VIP
        }

        $afterDiscount = $price - ($price * $discount);
        $finalPrice = $afterDiscount + ($afterDiscount * $tax);

        return round($finalPrice, 2);
    }
}
?>
