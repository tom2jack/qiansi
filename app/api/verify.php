<?php
/**
 * 验证系列接口
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-01-29
 * Time: 17:28
 */
namespace app\api;

use zhimiao\Config;
use zhimiao\Request;
use zhimiao\Utils;
use app\Service\Utils as APP_Utils;
use Gregwar\Captcha\CaptchaBuilder;
use Gregwar\Captcha\PhraseBuilder;
use AlibabaCloud\Client\AlibabaCloud;
use AlibabaCloud\Client\Exception\ClientException;
use AlibabaCloud\Client\Exception\ServerException;
use Respect\Validation\Validator as v;

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
        // $builder->setBackgroundColor(220, 210, 230);
        // $builder->setMaxAngle(8);
        // $builder->setMaxBehindLines(4);
        // $builder->setMaxFrontLines(4);
        $builder->setInterpolation(false);
        // 可以设置图片宽高及字体
        $builder->build();
        $result = [
            'image' => $builder->inline(),
            'token' => Utils::encrypt(json_encode(
                [
                    'code' => $code,
                    'key' => md5(Request::getIp()),
                    'ip' => Request::getIp(),
                    'expire' => time() + 1800,
                ]
            ))
        ];
        return [1, $result, "获取成功"];
    }

    /**
     * 获取手机验证码
     * @param string $phone
     * @param string $token
     * @param string $code
     * @return array|int
     */
    public function getSMSCode($phone = null, $token = null, $code = null)
    {
        if (!v::phone()->validate($phone)) {
            return [-5, null, '手机号无法识别'];
        }
        $VerifyService = new \app\Service\Verify();
        $check_result = $VerifyService->verifyImageCaptcha($token, $code);
        if (is_string($check_result)) {
            return [-5, null, $check_result];
        }
        // 增加过频锁
        $lock_id = 'verify/getSMSCode:'. $phone;
        if (!APP_Utils::cacheNumLock($lock_id, 1, 70)) {
            return [-5, null, '操作过频,手机号已锁定，请等待较长时间后重试'];
        }
        $phrase = new PhraseBuilder;
        // 设置验证码位数
        $code = $phrase->build(4);
        AlibabaCloud::accessKeyClient(Config::get('aliyun_sms.AccessKey'), Config::get('aliyun_sms.AccessSecret'))
            ->regionId(Config::get('aliyun_sms.RegionId'))
            ->asGlobalClient();
        try {
            $result = AlibabaCloud::rpcRequest()
                ->product('Dysmsapi')
                ->version('2017-05-25')
                ->action('SendSms')
                ->method('POST')
                ->options([
                    'query' => [
                        'RegionId' => Config::get('aliyun_sms.RegionId'),
                        'PhoneNumbers' => $phone,
                        'SignName' => Config::get('aliyun_sms.SignName'),
                        'TemplateCode' => Config::get('aliyun_sms.TemplateCode'),
                        'TemplateParam' => json_encode([
                            'verify' => $code
                        ]),
                    ],
                ])
                ->request()->toArray();
            if ($result['Code'] == 'OK') {
                $result = [
                    'token' => Utils::encrypt(json_encode(
                        [
                            'phone' => $phone,
                            'code' => $code,
                            'key' => md5(Request::getIp()),
                            'ip' => Request::getIp(),
                            'expire' => time() + 1800,
                        ]
                    ))
                ];
                return [1, $result, "获取成功"];
            }
        } catch (ClientException $e) {
            return [-7, null, $e->getErrorMessage()];
        } catch (ServerException $e) {
            return [-7, null, $e->getErrorMessage()];
        }
        return 0;
    }

}