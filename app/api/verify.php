<?php
/**
 * 验证系列接口
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-01-29
 * Time: 17:28
 */
namespace app\api;

use zhimiao\Utils;
use Gregwar\Captcha\CaptchaBuilder;
use Gregwar\Captcha\PhraseBuilder;

class verify {

    /**
     * 获取图像验证码
     * @return array
     */
    public function getImageCaptcha()
    {
        $phrase = new PhraseBuilder;
        // 设置验证码位数
        $code = $phrase->build(4);
        // 生成验证码图片的Builder对象，配置相应属性
        $builder = new CaptchaBuilder($code, $phrase);
        // 设置背景颜色
        $builder->setBackgroundColor(220, 210, 230);
        $builder->setMaxAngle(8);
        $builder->setMaxBehindLines(4);
        $builder->setMaxFrontLines(4);
        // 可以设置图片宽高及字体
        $builder->build($width = 100, $height = 40, $font = null);
        return [1, ['image' => $builder->inline(), "token" => Utils::encrypt($code)], "获取成功"];
    }

    public function getSMSCode($phone = null, $token = null, $code = null)
    {
        $token = Utils::decrypt($token);
        return $token;
    }
}