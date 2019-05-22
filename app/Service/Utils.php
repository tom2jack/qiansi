<?php
/**
 * 工具类
 * Created by PhpStorm.
 * User: 倒霉狐狸
 * Date: 2019-01-29
 * Time: 12:02
 */
namespace app\Service;
use \zhimiao\Config;
use zhimiao\Data;
use Respect\Validation\Validator as v;
class Utils {

    /**
     * 缓存锁 - 此方法需要配合计划任务解决死锁问题
     * @param string $key 键名
     * @param bool $begin true-上锁 false-解锁
     * @param int $expire 超时时间
     * @return bool 伪布尔
     */
    public static function cacheNumLock($key = '', $num = 5, $scale_time = 1800, $lock_time = 60)
    {
        $redis = Data::redis();
        $key = 'cacheNumLock:'. $key;
        // 如果计次为0，则清空锁
        if ($num === 0) {
            return (bool) $redis->del($key);
        }
        $data = $redis->get($key);
        $data = @json_decode($data, true);
        if (empty($data)) {
            $data = [
                'num' => 1,
                'last_time' => 0,
            ];
        }
        $out_num = $data['num'] - $num;
        if ($out_num > 0) {
            // 锁定时间次方增长，直至超过超时时间
            $lock_time = pow($lock_time, $out_num);
            if (time() < $data['last_time'] + $lock_time) {
                return false;
            }
            $scale_time = max([$scale_time, $lock_time]);
        }
        $data['num']++;
        $data['last_time'] = time();
        return (bool) $redis->set($key, json_encode($data), $scale_time);
    }

    /**
     * 创建加密的sessionKey
     * @param $uid
     * @param int $expire
     * @return string
     */
    public static function createSessionKey($uid, $expire = 86400)
    {
        $data = json_encode([
            'uid' => $uid,
            'expire' => time() + $expire
        ]);
        return \zhimiao\Utils::encrypt($data);
    }

    /**
     * 解析sessionKey
     * @param $sessionKey
     * @param bool $isExpire
     * @return mixed|string
     */
    public static function parseSessionKey($sessionKey, $isExpire = true)
    {
        $data = \zhimiao\Utils::decrypt($sessionKey);
        if (!v::json()->validate($data)) {
            return 'sessionKey无法解析';
        }
        $data = json_decode($data, true);
        if ($isExpire && time() > $data['expire']) {
            return 'sessionKey已经过期了';
        }
        return $data;
    }

    /**
     * 加密
     * @param $encryptStr 待加密内容
     * @param null $secret 密钥
     * @param string $iv 向量
     * @return string
     */
    static public function encrypt($encryptStr, $secret = null, $iv = '51ae84ba12c3a6ab991a89070555bae8')
    {
        $secret = $secret ?? Config::get('openssl_secret.key');
        $iv = hex2bin($iv);
        return openssl_encrypt($encryptStr, Config::get('openssl_secret.method'), $secret, 0, $iv);
    }

    /**
     * 解密
     * @param $encryptStr
     * @param null $secret
     * @param string $iv
     * @return string
     */
    static public function decrypt($encryptStr, $secret = null, $iv = '51ae84ba12c3a6ab991a89070555bae8')
    {
        $secret = $secret ?? Config::get('openssl_secret.key');
        $iv = hex2bin($iv);
        return openssl_decrypt($encryptStr, Config::get('openssl_secret.method'), $secret, 0, $iv);
    }

    /**
     * 获得随机字符串
     * @param $len 需要的长度
     * @param $special 是否需要特殊符号
     * @return string 返回随机字符串
     */
    static public function getRandomStr($len, $special = true){
        $chars = array(
            "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
            "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v",
            "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G",
            "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
            "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2",
            "3", "4", "5", "6", "7", "8", "9"
        );

        if($special) {
            $chars = array_merge($chars, array(
                "!", "@", "#", "$", "?", "|", "{", "/", ":", ";",
                "%", "^", "&", "*", "(", ")", "-", "_", "[", "]",
                "}", "<", ">", "~", "+", "=", ",", "."
            ));
        }

        $charsLen = count($chars) - 1;
        shuffle($chars);                            //打乱数组顺序
        $str = '';
        for($i=0; $i<$len; $i++){
            $str .= $chars[mt_rand(0, $charsLen)];    //随机取出一位
        }
        return $str;
    }
}